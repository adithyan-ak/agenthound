package rules

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// bundleOverridePath is the process-global override set by the CLI
// flag --rules-bundle. When non-empty, LoadFingerprints (the existing
// entry point every fingerprinter calls) merges bundle rules over the
// embedded set with same-id overrides winning.
//
// Set ONCE during CLI initialization (before any module's Fingerprint()
// call). Subsequent LoadFingerprints reads see the override transparently.
var (
	bundleOverridePathMu sync.RWMutex
	bundleOverridePath   string
)

// SetBundleOverridePath configures the process-global rules-bundle
// override. Pass an empty string to clear. Subsequent calls to
// LoadFingerprints merge bundle rules into the embedded set.
func SetBundleOverridePath(path string) {
	bundleOverridePathMu.Lock()
	defer bundleOverridePathMu.Unlock()
	bundleOverridePath = path
}

// getBundleOverridePath returns the current override or empty string.
// Internal to the package.
func getBundleOverridePath() string {
	bundleOverridePathMu.RLock()
	defer bundleOverridePathMu.RUnlock()
	return bundleOverridePath
}

// BundleSource identifies how a fingerprint rule entered the engine. The
// embedded set ships in the binary; "bundle" entries come from a
// --rules-bundle <path> override. The Source field on FingerprintRule
// already carries this for the embedded case ("builtin"); the override
// path replaces it with the absolute path of the bundle the rule came
// from so operators have a clear provenance trail.
const (
	BundleSourceBuiltin = "builtin"
)

// LoadFingerprintBundle reads fingerprint rules from a directory or
// tar.gz file and returns them. Same-id rules from the bundle WIN over
// embedded rules — operators can ship a rule fix without cutting a
// binary release. Used by the v0.3 --rules-bundle CLI flag.
//
// Cosign signature verification is the operator's responsibility for
// v0.3 — the loader does not validate signatures (mandatory verification
// lands in v0.5 once the release pipeline has cut at least one bundle).
// Operators should run cosign verify-blob against the tarball before
// pointing AgentHound at it; see docs/rules-bundle.md.
//
// path may be:
//   - a directory containing *.yaml files (each file = one rule)
//   - a .tar.gz file containing one *.yaml entry per rule
//
// Files that fail to parse are skipped with a warning logged via the
// caller's slog (this function does NOT log directly — it returns the
// successful subset and lets the caller decide policy).
func LoadFingerprintBundle(path string) ([]FingerprintRule, error) {
	if path == "" {
		return nil, errors.New("LoadFingerprintBundle: empty path")
	}
	st, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat bundle path: %w", err)
	}
	if st.IsDir() {
		return loadBundleFromDir(path)
	}
	return loadBundleFromTarball(path)
}

func loadBundleFromDir(dir string) ([]FingerprintRule, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read bundle dir: %w", err)
	}
	var rules []FingerprintRule
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		r, err := parseBundleRule(data, full)
		if err != nil {
			continue
		}
		rules = append(rules, *r)
	}
	return rules, nil
}

func loadBundleFromTarball(path string) ([]FingerprintRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open bundle tarball: %w", err)
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("gunzip bundle: %w", err)
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	var rules []FingerprintRule
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return rules, fmt.Errorf("tar read: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if !strings.HasSuffix(hdr.Name, ".yaml") {
			continue
		}
		// Cap per-file size at 1 MiB. Fingerprint YAMLs are tiny;
		// anything larger is suspicious.
		data, err := io.ReadAll(io.LimitReader(tr, 1<<20))
		if err != nil {
			continue
		}
		r, err := parseBundleRule(data, "bundle:"+hdr.Name)
		if err != nil {
			continue
		}
		rules = append(rules, *r)
	}
	return rules, nil
}

func parseBundleRule(data []byte, source string) (*FingerprintRule, error) {
	var r FingerprintRule
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("yaml unmarshal: %w", err)
	}
	r.Source = source
	if r.Version == 0 {
		r.Version = 2
	}
	return &r, nil
}

// MergeFingerprintRules merges a base rule set (typically the embedded
// builtin set) with an override set (from a --rules-bundle path). When
// the same ID appears in both, the override wins. Rules from base that
// don't appear in override pass through unchanged.
//
// This is the operational primitive that lets a hot-fix bundle replace
// a broken embedded rule without rebuilding the binary. The override
// semantics are explicit — same-id overrides land cleanly; new IDs
// from the bundle add to the set.
func MergeFingerprintRules(base, override []FingerprintRule) []FingerprintRule {
	byID := make(map[string]FingerprintRule, len(base)+len(override))
	for _, r := range base {
		byID[r.ID] = r
	}
	for _, r := range override {
		byID[r.ID] = r
	}
	out := make([]FingerprintRule, 0, len(byID))
	for _, r := range byID {
		out = append(out, r)
	}
	return out
}
