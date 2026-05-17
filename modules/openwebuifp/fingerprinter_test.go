package openwebuifp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

const owuiVersionBody = `{"version":"0.6.5"}`
const owuiConfigBodyWithBackend = `{"name":"Open WebUI","ollama":{"base_url":"http://ollama-backend:11434"}}`
const owuiConfigBodyNoBackend = `{"name":"Open WebUI","ollama":{}}`

func TestFingerprint_OpenWebUIHappyWithExposes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/version":
			_, _ = w.Write([]byte(owuiVersionBody))
		case "/api/config":
			_, _ = w.Write([]byte(owuiConfigBodyWithBackend))
		default:
			w.WriteHeader(404)
		}
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
	if res.ServiceKind != "openwebui" {
		t.Errorf("ServiceKind = %q, want openwebui", res.ServiceKind)
	}
	if res.IngestData == nil {
		t.Fatal("IngestData nil")
	}
	if len(res.IngestData.Graph.Nodes) != 2 {
		t.Fatalf("expected 2 nodes (OpenWebUI + placeholder Ollama), got %d", len(res.IngestData.Graph.Nodes))
	}
	if len(res.IngestData.Graph.Edges) != 1 {
		t.Fatalf("expected 1 EXPOSES edge, got %d", len(res.IngestData.Graph.Edges))
	}
	edge := res.IngestData.Graph.Edges[0]
	if edge.Kind != "EXPOSES" {
		t.Errorf("edge kind = %q, want EXPOSES", edge.Kind)
	}
	if edge.SourceKind != "OpenWebUIInstance" || edge.TargetKind != "OllamaInstance" {
		t.Errorf("edge endpoints = %s -> %s, want OpenWebUIInstance -> OllamaInstance", edge.SourceKind, edge.TargetKind)
	}
	if got, _ := edge.Properties["evidence"].(string); got != "http://ollama-backend:11434" {
		t.Errorf("edge evidence = %v, want http://ollama-backend:11434", edge.Properties["evidence"])
	}
}

func TestFingerprint_OpenWebUI_NoBackendStillMatches(t *testing.T) {
	// /api/config available but missing ollama.base_url — fingerprint
	// matches, no EXPOSES edge.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/version":
			_, _ = w.Write([]byte(owuiVersionBody))
		case "/api/config":
			_, _ = w.Write([]byte(owuiConfigBodyNoBackend))
		default:
			w.WriteHeader(404)
		}
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
	if len(res.IngestData.Graph.Nodes) != 1 {
		t.Errorf("expected 1 node (no placeholder Ollama), got %d", len(res.IngestData.Graph.Nodes))
	}
	if len(res.IngestData.Graph.Edges) != 0 {
		t.Errorf("expected 0 EXPOSES edges (no backend URL captured), got %d", len(res.IngestData.Graph.Edges))
	}
}

func TestFingerprint_OpenWebUI_ConfigLockedFallback(t *testing.T) {
	// /api/config 401 — fingerprinter must still match on /api/version alone.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/version":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(owuiVersionBody))
		case "/api/config":
			w.WriteHeader(401)
		default:
			w.WriteHeader(404)
		}
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
		t.Fatal("expected Matched=true even when /api/config is locked")
	}
	if len(res.IngestData.Graph.Edges) != 0 {
		t.Errorf("expected 0 EXPOSES edges when config locked, got %d", len(res.IngestData.Graph.Edges))
	}
}

func TestFingerprint_NotOpenWebUI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"object":"list","data":[]}`))
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
		t.Error("expected no match on non-OpenWebUI body")
	}
}

func TestCanonicalizeBackend(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"http://ollama:11434", "http://ollama:11434"},
		{"https://ollama.example.com", "https://ollama.example.com:11434"},
		{"ollama-backend:11434", "http://ollama-backend:11434"},
		{"ollama-backend", "http://ollama-backend:11434"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := canonicalizeBackend(tt.input)
			if got != tt.want {
				t.Errorf("canonicalizeBackend(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
