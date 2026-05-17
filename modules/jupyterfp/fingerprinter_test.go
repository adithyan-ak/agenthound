package jupyterfp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

const jupyterStatusBody = `{"started":"2026-04-01T12:00:00.000000Z","last_activity":"2026-04-01T12:34:56.000000Z","connections":3,"kernels":2}`

func TestFingerprint_JupyterHappy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/status" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(jupyterStatusBody))
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
	if res.ServiceKind != "jupyter" {
		t.Errorf("ServiceKind = %q, want jupyter", res.ServiceKind)
	}
	if res.AuthMethod != "token" {
		t.Errorf("AuthMethod = %q, want token", res.AuthMethod)
	}
	if res.IngestData == nil || len(res.IngestData.Graph.Nodes) != 1 {
		t.Fatalf("expected 1 ingest node, got %+v", res.IngestData)
	}
	node := res.IngestData.Graph.Nodes[0]
	if len(node.Kinds) != 2 || node.Kinds[0] != "JupyterServer" || node.Kinds[1] != "AIService" {
		t.Errorf("node.Kinds = %v, want [JupyterServer AIService]", node.Kinds)
	}
}

func TestFingerprint_NotJupyter(t *testing.T) {
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
		t.Error("expected no match — body is not Jupyter-shaped")
	}
}
