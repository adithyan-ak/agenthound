package appdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditStore struct {
	pool *pgxpool.Pool
}

type AuditEntry struct {
	ID        int64          `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	Action    string         `json:"action"`
	UserID    string         `json:"user_id,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}

type AuditFilter struct {
	Action string
	UserID string
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

func NewAuditStore(pool *pgxpool.Pool) *AuditStore {
	return &AuditStore{pool: pool}
}

func (s *AuditStore) Log(ctx context.Context, action, userID string, details map[string]any) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO audit_log (action, user_id, details) VALUES ($1, $2, $3)`,
		action, userID, details)
	if err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}

func (s *AuditStore) List(ctx context.Context, filter AuditFilter) ([]AuditEntry, error) {
	if filter.Limit <= 0 {
		filter.Limit = 100
	}

	query := `SELECT id, timestamp, action, COALESCE(user_id, ''), details FROM audit_log WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, filter.Action)
		argIdx++
	}
	if filter.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, filter.UserID)
		argIdx++
	}
	if filter.From != nil {
		query += fmt.Sprintf(" AND timestamp >= $%d", argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(" AND timestamp <= $%d", argIdx)
		args = append(args, *filter.To)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit log: %w", err)
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.Action, &e.UserID, &e.Details); err != nil {
			return nil, fmt.Errorf("scan audit row: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
