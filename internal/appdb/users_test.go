package appdb

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestIntegrationUsersCRUD(t *testing.T) {
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

	store := NewUserStore(pool)
	userID := "test-user-" + time.Now().Format("20060102150405")

	// Cleanup on exit
	defer func() {
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	}()

	// Create
	user := &model.User{
		ID:           userID,
		Username:     "testuser_" + time.Now().Format("150405"),
		PasswordHash: "$2a$10$fakehashfortest",
		Role:         "analyst",
		CreatedAt:    time.Now().UTC(),
	}
	if err := store.Create(ctx, user); err != nil {
		t.Fatalf("create: %v", err)
	}

	// GetByUsername
	got, err := store.GetByUsername(ctx, user.Username)
	if err != nil {
		t.Fatalf("get by username: %v", err)
	}
	if got.ID != userID {
		t.Errorf("id: got %q, want %q", got.ID, userID)
	}
	if got.Role != "analyst" {
		t.Errorf("role: got %q, want analyst", got.Role)
	}
	if got.PasswordHash != user.PasswordHash {
		t.Errorf("password_hash: got %q, want %q", got.PasswordHash, user.PasswordHash)
	}

	// GetByID
	got, err = store.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Username != user.Username {
		t.Errorf("username: got %q, want %q", got.Username, user.Username)
	}

	// UpdateLastLogin
	if err := store.UpdateLastLogin(ctx, userID); err != nil {
		t.Fatalf("update last login: %v", err)
	}
	got, err = store.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("get after login update: %v", err)
	}
	if got.LastLogin == nil {
		t.Error("last_login should not be nil after update")
	}

	// List
	users, err := store.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(users) == 0 {
		t.Error("expected at least 1 user in list")
	}

	// Count
	count, err := store.Count(ctx)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count < 1 {
		t.Errorf("count: got %d, want >= 1", count)
	}

	// Delete
	if err := store.Delete(ctx, userID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = store.GetByID(ctx, userID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestIntegrationUsersDuplicateUsername(t *testing.T) {
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

	store := NewUserStore(pool)
	ts := time.Now().Format("20060102150405")
	username := "dupuser_" + ts

	user1 := &model.User{
		ID:           "dup-test-1-" + ts,
		Username:     username,
		PasswordHash: "$2a$10$fakehash1",
		Role:         "analyst",
		CreatedAt:    time.Now().UTC(),
	}
	user2 := &model.User{
		ID:           "dup-test-2-" + ts,
		Username:     username,
		PasswordHash: "$2a$10$fakehash2",
		Role:         "viewer",
		CreatedAt:    time.Now().UTC(),
	}

	defer func() {
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE id IN ($1, $2)", user1.ID, user2.ID)
	}()

	if err := store.Create(ctx, user1); err != nil {
		t.Fatalf("create first user: %v", err)
	}

	err = store.Create(ctx, user2)
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
	if !strings.Contains(err.Error(), "create user") {
		t.Errorf("error should wrap with 'create user': %v", err)
	}
}
