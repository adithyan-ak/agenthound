package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestHealth_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/health" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	if err := c.Health(context.Background()); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestHealth_ConnectionRefused(t *testing.T) {
	c := New("http://127.0.0.1:1", "")
	err := c.Health(context.Background())
	if err == nil {
		t.Fatal("expected error for connection refused")
	}
	if got := err.Error(); got != "cannot reach server at http://127.0.0.1:1: is it running?" {
		t.Fatalf("unexpected error: %s", got)
	}
}

func TestIngest_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/ingest" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or wrong auth header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("missing content-type: %s", r.Header.Get("Content-Type"))
		}

		var data model.IngestData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if data.Meta.Collector != "mcp" {
			t.Errorf("expected collector 'mcp', got %q", data.Meta.Collector)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"scan_id":"test-123","nodes_written":10,"edges_written":20,"duration":1000000000}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "test-token")
	result, err := c.Ingest(context.Background(), &model.IngestData{
		Meta: model.IngestMeta{
			Version:   1,
			Type:      "agenthound-ingest",
			Collector: "mcp",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ScanID != "test-123" {
		t.Errorf("expected scan_id 'test-123', got %q", result.ScanID)
	}
	if result.NodesWritten != 10 {
		t.Errorf("expected 10 nodes, got %d", result.NodesWritten)
	}
	if result.EdgesWritten != 20 {
		t.Errorf("expected 20 edges, got %d", result.EdgesWritten)
	}
}

func TestIngest_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := New(srv.URL, "bad-token")
	_, err := c.Ingest(context.Background(), &model.IngestData{})
	if err == nil {
		t.Fatal("expected error on 401")
	}
	expected := "authentication failed: run 'agenthound setup' to reconfigure"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestIngest_RateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := New(srv.URL, "token")
	_, err := c.Ingest(context.Background(), &model.IngestData{})
	if err == nil {
		t.Fatal("expected error on 429")
	}
	expected := "rate limited by server, wait and retry"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGetFindings_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/analysis/findings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"abc123","severity":"critical","category":"Transitive Access","title":"Agent can reach resource","edge_kind":"CAN_REACH","source_name":"agent1","target_name":"resource1","confidence":1.0}]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "token")
	findings, err := c.GetFindings(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].ID != "abc123" {
		t.Errorf("expected id 'abc123', got %q", findings[0].ID)
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected severity 'critical', got %q", findings[0].Severity)
	}
	if findings[0].Confidence != 1.0 {
		t.Errorf("expected confidence 1.0, got %f", findings[0].Confidence)
	}
}

func TestGetFindings_WithSeverityFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("severity"); got != "critical" {
			t.Errorf("expected severity query param 'critical', got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "token")
	findings, err := c.GetFindings(context.Background(), "critical")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestLogin_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/auth/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("login should not send auth header when client has no token")
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Username != "admin" || req.Password != "secret" {
			t.Errorf("unexpected credentials: %s / %s", req.Username, req.Password)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"jwt-abc-123","expires_at":"2026-04-14T00:00:00Z","user":{"id":"u1","username":"admin","role":"admin"}}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	token, err := c.Login(context.Background(), "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "jwt-abc-123" {
		t.Errorf("expected token 'jwt-abc-123', got %q", token)
	}
}

func TestLogin_BadCredentials(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.Login(context.Background(), "admin", "wrong")
	if err == nil {
		t.Fatal("expected error on bad credentials")
	}
	expected := "invalid credentials"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestCreateToken_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/auth/tokens" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer jwt-token" {
			t.Errorf("wrong auth header: %s", r.Header.Get("Authorization"))
		}

		var req createTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "ci-pipeline" {
			t.Errorf("expected name 'ci-pipeline', got %q", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"token":"ah_abc123def456","id":"tok-1","name":"ci-pipeline"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "jwt-token")
	token, err := c.CreateToken(context.Background(), "ci-pipeline")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "ah_abc123def456" {
		t.Errorf("expected token 'ah_abc123def456', got %q", token)
	}
}

func TestGetPrebuilt_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/analysis/prebuilt/agents-shell-access" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("wrong auth header: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"agent":"claude-desktop","tool":"run_command","host":"localhost"}]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "test-token")
	rows, err := c.GetPrebuilt(context.Background(), "agents-shell-access")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["agent"] != "claude-desktop" {
		t.Errorf("expected agent 'claude-desktop', got %v", rows[0]["agent"])
	}
}
