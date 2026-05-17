package litellmloot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/common"
	"github.com/adithyan-ak/agenthound/sdk/ingest"
)

const fakeMasterKey = "sk-test-litellm-master-key-not-real"

// stubLiteLLM is a configurable test server simulating LiteLLM responses.
// Each handler controls what /model/info and /key/list return; tests
// mix-and-match to exercise happy path, partial failure, and lenient
// parsing paths.
type stubLiteLLM struct {
	t              *testing.T
	modelInfoBody  string
	modelInfoCode  int
	keyListBody    string
	keyListCode    int
	requireBearer  bool
	seenMethods    []string
	seenAuthHeader []string
}

func newStub(t *testing.T) *stubLiteLLM {
	return &stubLiteLLM{
		t:             t,
		modelInfoCode: 200,
		keyListCode:   200,
		requireBearer: true,
	}
}

func (s *stubLiteLLM) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.seenMethods = append(s.seenMethods, r.Method)
		s.seenAuthHeader = append(s.seenAuthHeader, r.Header.Get("Authorization"))
		if s.requireBearer && !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			w.WriteHeader(401)
			return
		}
		switch r.URL.Path {
		case "/model/info":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.modelInfoCode)
			_, _ = w.Write([]byte(s.modelInfoBody))
		case "/key/list":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.keyListCode)
			_, _ = w.Write([]byte(s.keyListBody))
		default:
			w.WriteHeader(404)
		}
	}
}

// happyPathBody is a representative LiteLLM /model/info response carrying
// three providers. Shape is deliberately simplified vs. the real LiteLLM
// payload — the looter parses leniently, so unknown fields are fine.
const happyPathModelInfo = `{
  "data": [
    {
      "model_name": "gpt-4",
      "litellm_params": {"model": "openai/gpt-4", "api_base": "https://api.openai.com/v1"},
      "model_info": {"litellm_provider": "openai"}
    },
    {
      "model_name": "claude-3",
      "litellm_params": {"model": "anthropic/claude-3-opus", "api_base": "https://api.anthropic.com"},
      "model_info": {"litellm_provider": "anthropic"}
    },
    {
      "model_name": "bedrock-claude",
      "litellm_params": {"model": "bedrock/anthropic.claude-v2"},
      "model_info": {"litellm_provider": "bedrock"}
    }
  ]
}`

const happyPathKeyList = `{
  "keys": [
    {"key_id": "vk-eng-team", "spend": 12.34, "models": ["gpt-4", "claude-3"]},
    {"key_id": "vk-data-team", "spend": 5.67, "models": ["claude-3"]}
  ]
}`

func TestLoot_HappyPath(t *testing.T) {
	s := newStub(t)
	s.modelInfoBody = happyPathModelInfo
	s.keyListBody = happyPathKeyList
	srv := httptest.NewServer(s.handler())
	defer srv.Close()

	l := &Looter{}
	res, err := l.Loot(context.Background(), action.Target{
		Kind:    "host",
		Address: strings.TrimPrefix(srv.URL, "http://"),
	}, action.LootOptions{
		Credentials:  map[string]string{"master_key": fakeMasterKey},
		EngagementID: "TEST-ENGAGEMENT",
		Timeout:      5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Loot: %v", err)
	}
	if res.IngestData == nil {
		t.Fatal("nil IngestData")
	}

	// Expect: 1 master + 3 upstream + 2 virtual = 6 Credential nodes,
	// 6 EXPOSES_CREDENTIAL edges.
	credCount := 0
	for _, n := range res.IngestData.Graph.Nodes {
		if len(n.Kinds) > 0 && n.Kinds[0] == "Credential" {
			credCount++
		}
	}
	if credCount != 6 {
		t.Errorf("Credential nodes = %d, want 6 (1 master + 3 upstream + 2 virtual)", credCount)
	}

	edgeCount := 0
	for _, e := range res.IngestData.Graph.Edges {
		if e.Kind == "EXPOSES_CREDENTIAL" {
			edgeCount++
		}
	}
	if edgeCount != 6 {
		t.Errorf("EXPOSES_CREDENTIAL edges = %d, want 6", edgeCount)
	}

	if len(res.PartialErrors) != 0 {
		t.Errorf("PartialErrors = %v, want []", res.PartialErrors)
	}
	if res.Summary.CredentialsFound != 6 {
		t.Errorf("Summary.CredentialsFound = %d, want 6", res.Summary.CredentialsFound)
	}

	// The master Credential MUST carry value_hash matching
	// HashCredentialValue(masterKey) — this is the cross-collector
	// merge primitive. Without this the credential-chain demo fails.
	var masterNode *ingest.Node
	for i := range res.IngestData.Graph.Nodes {
		if res.IngestData.Graph.Nodes[i].Properties["type"] == "master_key" {
			masterNode = &res.IngestData.Graph.Nodes[i]
			break
		}
	}
	if masterNode == nil {
		t.Fatal("master Credential node not emitted")
	}
	wantHash := common.HashCredentialValue(fakeMasterKey)
	if got := masterNode.Properties["value_hash"]; got != wantHash {
		t.Errorf("master value_hash = %v, want %v (cross-collector merge primitive)", got, wantHash)
	}

	// Master raw value must NOT be set when IncludeCredentialValues=false.
	if _, ok := masterNode.Properties["value"]; ok {
		t.Errorf("master node leaked raw value with IncludeCredentialValues=false")
	}

	// EngagementID must surface on edge evidence.
	for _, e := range res.IngestData.Graph.Edges {
		if e.Kind != "EXPOSES_CREDENTIAL" {
			continue
		}
		ev, _ := e.Properties["evidence"].(map[string]any)
		if ev["engagement_id"] != "TEST-ENGAGEMENT" {
			t.Errorf("edge evidence.engagement_id = %v, want TEST-ENGAGEMENT", ev["engagement_id"])
		}
	}
}

func TestLoot_PartialFailure_KeyListUnauthorized(t *testing.T) {
	s := newStub(t)
	s.modelInfoBody = happyPathModelInfo
	// /key/list 401 — common in production where the master key is
	// scoped down. The looter must record the error and continue,
	// returning the upstream credentials it did extract.
	s.keyListCode = 401
	s.keyListBody = `{"error": "unauthorized"}`
	srv := httptest.NewServer(s.handler())
	defer srv.Close()

	l := &Looter{}
	res, err := l.Loot(context.Background(), action.Target{
		Kind:    "host",
		Address: strings.TrimPrefix(srv.URL, "http://"),
	}, action.LootOptions{
		Credentials: map[string]string{"master_key": fakeMasterKey},
	})
	if err != nil {
		t.Fatalf("Loot returned err on partial failure: %v", err)
	}
	if len(res.PartialErrors) == 0 {
		t.Fatal("expected non-empty PartialErrors")
	}
	hasKeyListErr := false
	for _, pe := range res.PartialErrors {
		if strings.Contains(pe, "key/list") {
			hasKeyListErr = true
		}
	}
	if !hasKeyListErr {
		t.Errorf("PartialErrors missing key/list entry: %v", res.PartialErrors)
	}
	// Should still have master + 3 upstream credentials.
	credCount := 0
	for _, n := range res.IngestData.Graph.Nodes {
		if len(n.Kinds) > 0 && n.Kinds[0] == "Credential" {
			credCount++
		}
	}
	if credCount != 4 {
		t.Errorf("Credential nodes = %d, want 4 (1 master + 3 upstream)", credCount)
	}
}

func TestLoot_LenientModelInfoShape(t *testing.T) {
	// LiteLLM's /model/info shape has drifted across versions; the
	// looter must not fail-fast on unexpected fields. A response with
	// a single entry and no api_base / no api_key still produces an
	// upstream Credential with a synthetic value_hash.
	s := newStub(t)
	s.modelInfoBody = `{"data": [{"model_name": "gpt-4", "litellm_params": {"model": "openai/gpt-4"}, "model_info": {}}]}`
	s.keyListBody = `{"keys": []}`
	srv := httptest.NewServer(s.handler())
	defer srv.Close()

	l := &Looter{}
	res, err := l.Loot(context.Background(), action.Target{
		Kind:    "host",
		Address: strings.TrimPrefix(srv.URL, "http://"),
	}, action.LootOptions{
		Credentials: map[string]string{"master_key": fakeMasterKey},
	})
	if err != nil {
		t.Fatalf("Loot: %v", err)
	}
	credCount := 0
	for _, n := range res.IngestData.Graph.Nodes {
		if len(n.Kinds) > 0 && n.Kinds[0] == "Credential" {
			credCount++
		}
	}
	if credCount != 2 {
		t.Errorf("Credential nodes = %d, want 2 (master + 1 upstream)", credCount)
	}
}

func TestLoot_RequiresMasterKey(t *testing.T) {
	l := &Looter{}
	_, err := l.Loot(context.Background(), action.Target{Address: "127.0.0.1:4000"},
		action.LootOptions{})
	if err == nil {
		t.Fatal("expected error for missing master_key")
	}
	if !strings.Contains(err.Error(), "master") {
		t.Errorf("err = %v, want one mentioning 'master'", err)
	}
}

func TestLoot_IncludeCredentialValues(t *testing.T) {
	// When the operator opts in, the master Credential carries the raw
	// value too. The merge-primitive value_hash is unchanged.
	s := newStub(t)
	s.modelInfoBody = `{"data": []}`
	s.keyListBody = `{"keys": []}`
	srv := httptest.NewServer(s.handler())
	defer srv.Close()

	l := &Looter{}
	res, err := l.Loot(context.Background(), action.Target{
		Address: strings.TrimPrefix(srv.URL, "http://"),
	}, action.LootOptions{
		Credentials:             map[string]string{"master_key": fakeMasterKey},
		IncludeCredentialValues: true,
	})
	if err != nil {
		t.Fatalf("Loot: %v", err)
	}
	var master *ingest.Node
	for i := range res.IngestData.Graph.Nodes {
		if res.IngestData.Graph.Nodes[i].Properties["type"] == "master_key" {
			master = &res.IngestData.Graph.Nodes[i]
		}
	}
	if master == nil {
		t.Fatal("master node missing")
	}
	if got := master.Properties["value"]; got != fakeMasterKey {
		t.Errorf("master.value = %v, want raw key with IncludeCredentialValues=true", got)
	}
}
