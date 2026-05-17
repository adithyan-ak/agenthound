// Package litellmfp implements the v0.2 LiteLLM fingerprinter module.
// It probes a Target's Address (expected shape "host:4000") with
// GET /health/liveliness and emits a multi-label
// :LiteLLMGateway:AIService node when the response matches LiteLLM's
// canonical liveliness body ("I'm alive!").
//
// The probe semantics, matchers, and node properties live in
// sdk/rules/builtin/fingerprints/litellm.yaml. This package is the
// dispatcher that loads that YAML, locates it by service_kind, and
// runs it via sdk/rules.RunFingerprint.
//
// LiteLLM is the v0.2 Looter target — once a LiteLLMGateway is in the
// graph, the operator runs `agenthound loot --type litellm` against it
// to extract upstream provider keys (Phase 4).
package litellmfp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/ingest"
	"github.com/adithyan-ak/agenthound/sdk/rules"
)

const (
	DefaultPort         = 4000
	DefaultProbeTimeout = 5 * time.Second
)

type Fingerprinter struct {
	rule *rules.FingerprintRule
}

// New loads the LiteLLM fingerprint rule and returns a ready-to-use
// Fingerprinter. Same shape as Ollama's New(); the duplication is
// intentional in v0.2 — refactor when a third fingerprinter ships.
func New() (*Fingerprinter, error) {
	all, err := rules.LoadFingerprints()
	if err != nil {
		return nil, fmt.Errorf("load fingerprint rules: %w", err)
	}
	for _, r := range all {
		if r.ServiceKind == "litellm" {
			rule := r
			if errs := rules.ValidateFingerprint(rule); len(errs) > 0 {
				return nil, fmt.Errorf("litellm rule invalid: %v", errs)
			}
			return &Fingerprinter{rule: &rule}, nil
		}
	}
	return nil, errors.New("litellm fingerprint rule not found in builtin set")
}

// Fingerprint runs the LiteLLM probe against t.Address. See
// modules/ollamafp for the design rationale; this is the same shape
// against a different rule.
func (f *Fingerprinter) Fingerprint(ctx context.Context, t action.Target) (*action.FingerprintResult, error) {
	if f.rule == nil {
		return nil, errors.New("litellm fingerprinter: rule not loaded")
	}
	host, port := splitHostPort(t.Address, DefaultPort)
	scheme := "http"
	if s, ok := t.Meta["scheme"]; ok && s != "" {
		scheme = s
	}
	baseURL := fmt.Sprintf("%s://%s:%d", scheme, host, port)

	client := rules.DefaultFingerprintHTTPClient(DefaultProbeTimeout)
	res, err := rules.RunFingerprint(ctx, client, baseURL, *f.rule)
	if err != nil {
		slog.Debug("litellm fingerprint probe error",
			"target", t.Address, "error", err)
		return &action.FingerprintResult{Matched: false}, nil
	}
	if !res.Matched {
		return &action.FingerprintResult{Matched: false}, nil
	}

	endpoint := baseURL
	objectID := ingest.ComputeNodeID("LiteLLMGateway", endpoint)

	props := map[string]any{
		"objectid":       objectID,
		"name":           host,
		"endpoint":       endpoint,
		"discovered_via": "network_scan",
		"service_kind":   "litellm",
	}
	for k, v := range res.Properties {
		props[k] = v
	}
	if _, ok := props["auth_method"]; !ok {
		props["auth_method"] = "master_key"
	}

	node := ingest.Node{
		ID:         objectID,
		Kinds:      append([]string{}, res.NodeKinds...),
		Properties: props,
	}
	out := &ingest.IngestData{
		Graph: ingest.GraphData{Nodes: []ingest.Node{node}},
	}

	return &action.FingerprintResult{
		Matched:     true,
		ServiceKind: "litellm",
		Version:     res.Properties["version"],
		AuthMethod:  res.Properties["auth_method"],
		IngestData:  out,
		Properties:  res.Properties,
	}, nil
}

// splitHostPort is duplicated from modules/ollamafp; refactor to a
// shared helper when fingerprinter #3 ships.
func splitHostPort(addr string, defaultPort int) (string, int) {
	addr = strings.TrimSpace(addr)
	if strings.Contains(addr, "://") {
		if u, err := url.Parse(addr); err == nil && u.Host != "" {
			return splitHostPort(u.Host, defaultPort)
		}
	}
	if i := strings.LastIndexByte(addr, ':'); i > 0 {
		host := addr[:i]
		var p int
		_, _ = fmt.Sscanf(addr[i+1:], "%d", &p)
		if p > 0 {
			host = strings.TrimPrefix(strings.TrimSuffix(host, "]"), "[")
			return host, p
		}
	}
	host := strings.TrimPrefix(strings.TrimSuffix(addr, "]"), "[")
	return host, defaultPort
}

var _ action.Fingerprinter = (*Fingerprinter)(nil)
