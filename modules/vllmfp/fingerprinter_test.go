package vllmfp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

const vllmModelsBody = `{"object":"list","data":[{"id":"meta-llama/Llama-3.1-8B","object":"model","created":1700000000,"owned_by":"vllm"}]}`

func TestFingerprint_VLLMHappy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/v1/models" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vllmModelsBody))
			return
		}
		w.WriteHeader(404)
	}))
	defer srv.Close()

	f, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	res, err := f.Fingerprint(context.Background(), action.Target{
		Kind:    "host",
		Address: strings.TrimPrefix(srv.URL, "http://"),
	})
	if err != nil {
		t.Fatalf("Fingerprint: %v", err)
	}
	if !res.Matched {
		t.Fatal("expected Matched=true")
	}
	if res.ServiceKind != "vllm" {
		t.Errorf("ServiceKind = %q, want vllm", res.ServiceKind)
	}
	if res.AuthMethod != "none" {
		t.Errorf("AuthMethod = %q, want none", res.AuthMethod)
	}
	if res.IngestData == nil || len(res.IngestData.Graph.Nodes) != 1 {
		t.Fatalf("expected 1 ingest node, got %+v", res.IngestData)
	}
	node := res.IngestData.Graph.Nodes[0]
	if len(node.Kinds) != 2 || node.Kinds[0] != "VLLMInstance" || node.Kinds[1] != "AIService" {
		t.Errorf("node.Kinds = %v, want [VLLMInstance AIService]", node.Kinds)
	}
	if got := node.Properties["service_kind"]; got != "vllm" {
		t.Errorf("properties.service_kind = %v, want vllm", got)
	}
	if !strings.HasPrefix(node.ID, "sha256:") {
		t.Errorf("node.ID = %q, want sha256: prefix", node.ID)
	}
}

func TestFingerprint_NotVLLM(t *testing.T) {
	// /v1/models returns Ollama-style JSON (no "object: list").
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"version":"0.5.1"}`))
	}))
	defer srv.Close()

	f, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	res, err := f.Fingerprint(context.Background(), action.Target{
		Kind:    "host",
		Address: strings.TrimPrefix(srv.URL, "http://"),
	})
	if err != nil {
		t.Fatalf("Fingerprint err = %v", err)
	}
	if res.Matched {
		t.Error("expected no match — body is Ollama-shaped, not vLLM")
	}
}

func TestFingerprint_NetworkError(t *testing.T) {
	f, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	res, err := f.Fingerprint(context.Background(), action.Target{
		Kind:    "host",
		Address: "127.0.0.1:1",
	})
	if err != nil {
		t.Fatalf("Fingerprint should not error on a closed port; got %v", err)
	}
	if res.Matched {
		t.Error("expected no match on closed port")
	}
}
