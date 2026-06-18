package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adithyan-ak/agenthound/server/model"
)

func TestHandleCreateScan_MissingCollector(t *testing.T) {
	h := NewScanHandler(nil, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/scans", []byte(`{}`))
	h.HandleCreate(w, r)

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

func TestHandleGetScan_EmptyID(t *testing.T) {
	h := NewScanHandler(nil, nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodGet, "/api/v1/scans/", nil)
	r = withChiURLParam(r, "id", "")
	h.HandleGet(w, r)

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

func TestHandleDeleteScan_GraphCleanupFailureDoesNotDeleteScan(t *testing.T) {
	store := &fakeScanStoreForHandler{scan: &model.Scan{ID: "scan-1", Collector: "mcp", Status: model.ScanStatusCompleted}}
	h := &ScanHandler{
		scanStore: store,
		graphDB:   &mockGraphDB{writeErr: errors.New("neo4j down")},
	}
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodDelete, "/api/v1/scans/scan-1", nil)
	r = withChiURLParam(r, "id", "scan-1")

	h.HandleDelete(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if store.deleted {
		t.Fatal("scan store DeleteScan should not be called when graph cleanup fails")
	}
}

type fakeScanStoreForHandler struct {
	scan    *model.Scan
	deleted bool
}

func (s *fakeScanStoreForHandler) ListScans(_ context.Context, _, _ int) ([]model.Scan, error) {
	if s.scan == nil {
		return nil, nil
	}
	return []model.Scan{*s.scan}, nil
}

func (s *fakeScanStoreForHandler) GetScan(_ context.Context, _ string) (*model.Scan, error) {
	if s.scan == nil {
		return nil, errors.New("not found")
	}
	return s.scan, nil
}

func (s *fakeScanStoreForHandler) CreateScan(_ context.Context, scan *model.Scan) error {
	s.scan = scan
	return nil
}

func (s *fakeScanStoreForHandler) DeleteScan(_ context.Context, _ string) error {
	s.deleted = true
	return nil
}
