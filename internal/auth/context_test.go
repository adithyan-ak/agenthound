package auth

import (
	"context"
	"testing"
	"time"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestContextWithUser_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	user := &model.User{
		ID:           "user-001",
		Username:     "admin",
		PasswordHash: "$2a$10$fake",
		Role:         "admin",
		CreatedAt:    now,
	}

	ctx := ContextWithUser(context.Background(), user)
	got := UserFromContext(ctx)

	if got == nil {
		t.Fatal("UserFromContext returned nil")
	}
	if got.ID != user.ID {
		t.Errorf("ID = %q, want %q", got.ID, user.ID)
	}
	if got.Username != user.Username {
		t.Errorf("Username = %q, want %q", got.Username, user.Username)
	}
	if got.Role != user.Role {
		t.Errorf("Role = %q, want %q", got.Role, user.Role)
	}
	if !got.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, user.CreatedAt)
	}
	if got != user {
		t.Error("expected same pointer identity")
	}
}

func TestUserFromContext_Empty(t *testing.T) {
	got := UserFromContext(context.Background())
	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}
