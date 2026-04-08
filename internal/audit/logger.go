package audit

import (
	"context"

	"github.com/adithyan-ak/agenthound/internal/appdb"
	"github.com/adithyan-ak/agenthound/internal/auth"
)

type Logger struct {
	store *appdb.AuditStore
}

func NewLogger(store *appdb.AuditStore) *Logger {
	return &Logger{store: store}
}

func (l *Logger) Log(ctx context.Context, action string, details map[string]any) error {
	var userID string
	if u := auth.UserFromContext(ctx); u != nil {
		userID = u.ID
	}
	return l.store.Log(ctx, action, userID, details)
}
