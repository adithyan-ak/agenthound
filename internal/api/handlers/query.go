package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/adithyan-ak/agenthound/internal/audit"
	"github.com/adithyan-ak/agenthound/internal/graph"
)

type QueryHandler struct {
	reader *graph.Reader
	audit  *audit.Logger
}

func NewQueryHandler(reader *graph.Reader, auditLog *audit.Logger) *QueryHandler {
	return &QueryHandler{reader: reader, audit: auditLog}
}

type queryRequest struct {
	Cypher string         `json:"cypher"`
	Params map[string]any `json:"params"`
}

func (h *QueryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req queryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "invalid JSON payload")
		return
	}
	if req.Cypher == "" {
		WriteValidationError(w, "cypher query is required")
		return
	}

	if h.audit != nil {
		cypher := req.Cypher
		if len(cypher) > 500 {
			cypher = cypher[:500]
		}
		if err := h.audit.Log(r.Context(), "query.execute", map[string]any{
			"cypher": cypher,
		}); err != nil {
			slog.Warn("audit log failed", "error", err)
		}
	}

	rows, err := h.reader.Query(r.Context(), req.Cypher, req.Params)
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{"rows": rows})
}
