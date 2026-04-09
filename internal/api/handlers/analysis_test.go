package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandleShortestPath_MissingSource(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{}, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/analysis/shortest-path", []byte(`{}`))
	h.HandleShortestPath(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}

func TestHandleShortestPath_InvalidKind(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{}, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/analysis/shortest-path",
		[]byte(`{"source":"x","source_kind":"INVALID"}`))
	h.HandleShortestPath(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleFindings_Empty(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{queryResult: nil}, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodGet, "/api/v1/analysis/findings", nil)
	h.HandleFindings(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var findings []any
	if err := json.NewDecoder(w.Body).Decode(&findings); err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestHandleListPreBuilt(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{}, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodGet, "/api/v1/analysis/prebuilt", nil)
	h.HandleListPreBuilt(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var queries []any
	if err := json.NewDecoder(w.Body).Decode(&queries); err != nil {
		t.Fatal(err)
	}
	if len(queries) != 17 {
		t.Fatalf("expected 17 pre-built queries, got %d", len(queries))
	}
}

func TestHandlePreBuilt_NotFound(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{}, nil)
	router := chi.NewRouter()
	router.Get("/api/v1/analysis/prebuilt/{id}", h.HandlePreBuilt)

	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodGet, "/api/v1/analysis/prebuilt/nonexistent-query", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %s", resp.Error.Code)
	}
}

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		targetKind string
		wantKind   string
		wantName   string
	}{
		{name: "both empty", target: "", targetKind: "", wantKind: "", wantName: ""},
		{name: "colon-separated target", target: "MCPServer:myserver", targetKind: "", wantKind: "MCPServer", wantName: "myserver"},
		{name: "plain target with kind", target: "myserver", targetKind: "MCPServer", wantKind: "MCPServer", wantName: "myserver"},
		{name: "empty target with kind", target: "", targetKind: "MCPServer", wantKind: "MCPServer", wantName: ""},
		{name: "colon target ignored when kind set", target: "MCPServer:myserver", targetKind: "A2AAgent", wantKind: "A2AAgent", wantName: "MCPServer:myserver"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKind, gotName := parseTarget(tt.target, tt.targetKind)
			if gotKind != tt.wantKind || gotName != tt.wantName {
				t.Errorf("parseTarget(%q, %q) = (%q, %q), want (%q, %q)",
					tt.target, tt.targetKind, gotKind, gotName, tt.wantKind, tt.wantName)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name       string
		val        int
		min        int
		max        int
		defaultVal int
		want       int
	}{
		{name: "zero returns default", val: 0, min: 1, max: 20, defaultVal: 10, want: 10},
		{name: "in range returns val", val: 5, min: 1, max: 20, defaultVal: 10, want: 5},
		{name: "negative returns default", val: -3, min: 1, max: 20, defaultVal: 10, want: 10},
		{name: "exceeds max clamped", val: 50, min: 1, max: 20, defaultVal: 10, want: 20},
		{name: "below min clamped", val: 1, min: 5, max: 20, defaultVal: 10, want: 5},
		{name: "exactly min", val: 1, min: 1, max: 20, defaultVal: 10, want: 1},
		{name: "exactly max", val: 20, min: 1, max: 20, defaultVal: 10, want: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clamp(tt.val, tt.min, tt.max, tt.defaultVal)
			if got != tt.want {
				t.Errorf("clamp(%d, %d, %d, %d) = %d, want %d",
					tt.val, tt.min, tt.max, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestHandleAllPaths_MissingSource(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{}, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/analysis/all-paths", []byte(`{}`))
	h.HandleAllPaths(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}

func TestHandleWeightedPath_MissingFields(t *testing.T) {
	h := NewAnalysisHandler(&mockGraphDB{}, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/analysis/weighted-path", []byte(`{}`))
	h.HandleWeightedPath(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}
