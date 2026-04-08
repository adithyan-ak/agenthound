package appdb

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestIntegrationAuditLog(t *testing.T) {
	skipIfNoPG(t)
	ctx := context.Background()

	pool, err := NewPool(os.Getenv("AGENTHOUND_PG_URI"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := RunMigrations(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	store := NewAuditStore(pool)
	tag := time.Now().Format("20060102150405")

	// Log entry
	details := map[string]any{"target": "scan-123"}
	if err := store.Log(ctx, "ingest.started", "user-"+tag, details); err != nil {
		t.Fatalf("log: %v", err)
	}

	// Log a second entry with different action
	if err := store.Log(ctx, "scan.completed", "user-"+tag, nil); err != nil {
		t.Fatalf("log second: %v", err)
	}

	// List all (unfiltered, scoped by user to avoid cross-test interference)
	entries, err := store.List(ctx, AuditFilter{UserID: "user-" + tag})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("list all: got %d entries, want 2", len(entries))
	}

	// List filtered by action
	entries, err = store.List(ctx, AuditFilter{Action: "ingest.started", UserID: "user-" + tag})
	if err != nil {
		t.Fatalf("list by action: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("list by action: got %d entries, want 1", len(entries))
	}
	if entries[0].Action != "ingest.started" {
		t.Errorf("action: got %q, want ingest.started", entries[0].Action)
	}
	if entries[0].UserID != "user-"+tag {
		t.Errorf("user_id: got %q, want user-%s", entries[0].UserID, tag)
	}

	// List filtered by time range
	before := time.Now().UTC().Add(-1 * time.Minute)
	after := time.Now().UTC().Add(1 * time.Minute)
	entries, err = store.List(ctx, AuditFilter{
		UserID: "user-" + tag,
		From:   &before,
		To:     &after,
	})
	if err != nil {
		t.Fatalf("list by time: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("list by time: got %d entries, want 2", len(entries))
	}

	// List with time range that excludes entries
	future := time.Now().UTC().Add(1 * time.Hour)
	farFuture := future.Add(1 * time.Hour)
	entries, err = store.List(ctx, AuditFilter{
		UserID: "user-" + tag,
		From:   &future,
		To:     &farFuture,
	})
	if err != nil {
		t.Fatalf("list empty range: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("list empty range: got %d entries, want 0", len(entries))
	}

	// Cleanup
	_, _ = pool.Exec(ctx, "DELETE FROM audit_log WHERE user_id = $1", "user-"+tag)
}
