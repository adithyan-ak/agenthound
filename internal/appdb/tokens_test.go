package appdb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestIntegrationTokenCreate(t *testing.T) {
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

	userID := "test-user-tok-create-" + time.Now().Format("20060102150405.000")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role) VALUES ($1, $2, $3, $4)`,
		userID, "tokuser-"+userID, "fakehash", "analyst")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "DELETE FROM api_tokens WHERE user_id = $1", userID)
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	store := NewTokenStore(pool)

	token := &model.APIToken{
		ID:        "tok-create-" + time.Now().Format("20060102150405.000"),
		UserID:    userID,
		TokenHash: "sha256-create-test",
		Name:      "test-token",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.Create(ctx, token); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.GetByHash(ctx, "sha256-create-test")
	if err != nil {
		t.Fatalf("get by hash: %v", err)
	}
	if got.ID != token.ID {
		t.Errorf("id: got %q, want %q", got.ID, token.ID)
	}
	if got.UserID != userID {
		t.Errorf("user_id: got %q, want %q", got.UserID, userID)
	}
	if got.Name != "test-token" {
		t.Errorf("name: got %q, want %q", got.Name, "test-token")
	}
	if got.ExpiresAt != nil {
		t.Errorf("expires_at: expected nil, got %v", got.ExpiresAt)
	}
}

func TestIntegrationTokenCreateWithExpiry(t *testing.T) {
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

	userID := "test-user-tok-expiry-" + time.Now().Format("20060102150405.000")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role) VALUES ($1, $2, $3, $4)`,
		userID, "tokuser-"+userID, "fakehash", "analyst")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "DELETE FROM api_tokens WHERE user_id = $1", userID)
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	store := NewTokenStore(pool)

	expiry := time.Now().UTC().Add(24 * time.Hour)
	token := &model.APIToken{
		ID:        "tok-expiry-" + time.Now().Format("20060102150405.000"),
		UserID:    userID,
		TokenHash: "sha256-expiry-test",
		Name:      "expiring-token",
		CreatedAt: time.Now().UTC(),
		ExpiresAt: &expiry,
	}
	if err := store.Create(ctx, token); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.GetByHash(ctx, "sha256-expiry-test")
	if err != nil {
		t.Fatalf("get by hash: %v", err)
	}
	if got.ExpiresAt == nil {
		t.Fatal("expires_at: expected non-nil")
	}
	if got.ExpiresAt.Sub(expiry).Abs() > time.Second {
		t.Errorf("expires_at: got %v, want ~%v", got.ExpiresAt, expiry)
	}
}

func TestIntegrationTokenGetByHashNotFound(t *testing.T) {
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

	store := NewTokenStore(pool)

	_, err = store.GetByHash(ctx, "nonexistent-hash")
	if err == nil {
		t.Fatal("expected error for nonexistent hash, got nil")
	}
}

func TestIntegrationTokenUpdateLastUsed(t *testing.T) {
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

	userID := "test-user-tok-lastused-" + time.Now().Format("20060102150405.000")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role) VALUES ($1, $2, $3, $4)`,
		userID, "tokuser-"+userID, "fakehash", "analyst")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "DELETE FROM api_tokens WHERE user_id = $1", userID)
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	store := NewTokenStore(pool)

	token := &model.APIToken{
		ID:        "tok-lastused-" + time.Now().Format("20060102150405.000"),
		UserID:    userID,
		TokenHash: "sha256-lastused-test",
		Name:      "lastused-token",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.Create(ctx, token); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.GetByHash(ctx, "sha256-lastused-test")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.LastUsed != nil {
		t.Errorf("last_used: expected nil before update, got %v", got.LastUsed)
	}

	if err := store.UpdateLastUsed(ctx, token.ID); err != nil {
		t.Fatalf("update last_used: %v", err)
	}

	got, err = store.GetByHash(ctx, "sha256-lastused-test")
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if got.LastUsed == nil {
		t.Fatal("last_used: expected non-nil after update")
	}
	if time.Since(*got.LastUsed) > 5*time.Second {
		t.Errorf("last_used: expected recent timestamp, got %v", got.LastUsed)
	}
}

func TestIntegrationTokenListByUser(t *testing.T) {
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

	userID := "test-user-tok-list-" + time.Now().Format("20060102150405.000")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role) VALUES ($1, $2, $3, $4)`,
		userID, "tokuser-"+userID, "fakehash", "analyst")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "DELETE FROM api_tokens WHERE user_id = $1", userID)
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	store := NewTokenStore(pool)

	// Empty list
	tokens, err := store.ListByUser(ctx, userID)
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}

	// Create two tokens
	now := time.Now().UTC()
	for i, name := range []string{"first-token", "second-token"} {
		tok := &model.APIToken{
			ID:        "tok-list-" + time.Now().Format("20060102150405.000") + "-" + name,
			UserID:    userID,
			TokenHash: "sha256-list-" + name,
			Name:      name,
			CreatedAt: now.Add(time.Duration(i) * time.Second),
		}
		if err := store.Create(ctx, tok); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}

	tokens, err = store.ListByUser(ctx, userID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	// ORDER BY created_at DESC
	if tokens[0].Name != "second-token" {
		t.Errorf("first result: got %q, want second-token", tokens[0].Name)
	}
	if tokens[1].Name != "first-token" {
		t.Errorf("second result: got %q, want first-token", tokens[1].Name)
	}
}

func TestIntegrationTokenDelete(t *testing.T) {
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

	userID := "test-user-tok-delete-" + time.Now().Format("20060102150405.000")
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role) VALUES ($1, $2, $3, $4)`,
		userID, "tokuser-"+userID, "fakehash", "analyst")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "DELETE FROM api_tokens WHERE user_id = $1", userID)
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	store := NewTokenStore(pool)

	token := &model.APIToken{
		ID:        "tok-delete-" + time.Now().Format("20060102150405.000"),
		UserID:    userID,
		TokenHash: "sha256-delete-test",
		Name:      "delete-me",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.Create(ctx, token); err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := store.Delete(ctx, token.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = store.GetByHash(ctx, "sha256-delete-test")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}
