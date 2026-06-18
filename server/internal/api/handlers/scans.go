package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adithyan-ak/agenthound/server/internal/appdb"
	"github.com/adithyan-ak/agenthound/server/internal/graph"
	"github.com/adithyan-ak/agenthound/server/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ScanHandler struct {
	scanStore scanStore
	graphDB   graph.GraphDB
}

type scanStore interface {
	ListScans(ctx context.Context, limit, offset int) ([]model.Scan, error)
	GetScan(ctx context.Context, id string) (*model.Scan, error)
	CreateScan(ctx context.Context, scan *model.Scan) error
	DeleteScan(ctx context.Context, id string) error
}

func NewScanHandler(store *appdb.ScanStore, graphDB graph.GraphDB) *ScanHandler {
	return &ScanHandler{scanStore: store, graphDB: graphDB}
}

func (h *ScanHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 50)
	offset := parseIntParam(r, "offset", 0)

	scans, err := h.scanStore.ListScans(r.Context(), limit, offset)
	if err != nil {
		WriteInternalError(w, r, fmt.Errorf("list scans: %w", err))
		return
	}
	if scans == nil {
		scans = []model.Scan{}
	}
	WriteJSON(w, http.StatusOK, scans)
}

func (h *ScanHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		WriteValidationError(w, "scan id is required")
		return
	}

	scan, err := h.scanStore.GetScan(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WriteNotFound(w, "scan not found")
			return
		}
		WriteInternalError(w, r, fmt.Errorf("get scan: %w", err))
		return
	}
	WriteJSON(w, http.StatusOK, scan)
}

type createScanRequest struct {
	Collector string         `json:"collector"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

func (h *ScanHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req createScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "invalid request body")
		return
	}
	if req.Collector == "" {
		WriteValidationError(w, "collector is required")
		return
	}
	validCollectors := map[string]bool{"mcp": true, "a2a": true, "config": true}
	if !validCollectors[req.Collector] {
		WriteValidationError(w, "collector must be one of: mcp, a2a, config")
		return
	}

	scan := model.Scan{
		ID:        uuid.New().String(),
		Collector: req.Collector,
		Status:    model.ScanStatusPending,
		StartedAt: time.Now().UTC(),
		Metadata:  req.Metadata,
	}

	if err := h.scanStore.CreateScan(r.Context(), &scan); err != nil {
		WriteInternalError(w, r, fmt.Errorf("create scan: %w", err))
		return
	}

	WriteJSON(w, http.StatusCreated, scan)
}

func (h *ScanHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		WriteValidationError(w, "scan id is required")
		return
	}

	// Verify the scan exists before doing graph cleanup.
	if _, err := h.scanStore.GetScan(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WriteNotFound(w, "scan not found")
			return
		}
		WriteInternalError(w, r, fmt.Errorf("get scan: %w", err))
		return
	}

	if err := h.deleteScanGraphData(r.Context(), id); err != nil {
		WriteInternalError(w, r, err)
		return
	}

	// Delete PG scan record.
	if err := h.scanStore.DeleteScan(r.Context(), id); err != nil {
		WriteInternalError(w, r, fmt.Errorf("delete scan: %w", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ScanHandler) deleteScanGraphData(ctx context.Context, id string) error {
	if h.graphDB == nil {
		return fmt.Errorf("delete scan graph data: graph database unavailable")
	}
	if _, err := h.graphDB.ExecuteWrite(ctx,
		`MATCH ()-[r]->() WHERE r.scan_id = $scan_id DELETE r RETURN count(r) AS deleted`,
		map[string]any{"scan_id": id}); err != nil {
		return fmt.Errorf("delete scan graph edges: %w", err)
	}
	if _, err := h.graphDB.ExecuteWrite(ctx,
		`MATCH (n) WHERE n.scan_id = $scan_id
		 AND NOT EXISTS { MATCH (n)-[]-() }
		 DELETE n RETURN count(n) AS deleted`,
		map[string]any{"scan_id": id}); err != nil {
		return fmt.Errorf("delete scan graph nodes: %w", err)
	}
	return nil
}
