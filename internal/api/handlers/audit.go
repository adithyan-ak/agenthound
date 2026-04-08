package handlers

import (
	"net/http"
	"time"

	"github.com/adithyan-ak/agenthound/internal/appdb"
)

type AuditHandler struct {
	store *appdb.AuditStore
}

func NewAuditHandler(store *appdb.AuditStore) *AuditHandler {
	return &AuditHandler{store: store}
}

func (h *AuditHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := appdb.AuditFilter{
		Action: q.Get("action"),
		UserID: q.Get("user_id"),
		Limit:  parseIntParam(r, "limit", 100),
		Offset: parseIntParam(r, "offset", 0),
	}

	if v := q.Get("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			WriteValidationError(w, "invalid 'from' timestamp: must be RFC3339")
			return
		}
		filter.From = &t
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			WriteValidationError(w, "invalid 'to' timestamp: must be RFC3339")
			return
		}
		filter.To = &t
	}

	entries, err := h.store.List(r.Context(), filter)
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}
	if entries == nil {
		entries = []appdb.AuditEntry{}
	}
	WriteJSON(w, http.StatusOK, entries)
}
