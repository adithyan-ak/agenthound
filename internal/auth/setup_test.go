package auth

import (
	"context"
	"os"
	"testing"

	"github.com/adithyan-ak/agenthound/internal/appdb"
)

func skipIfNoPG(t *testing.T) {
	t.Helper()
	if os.Getenv("AGENTHOUND_PG_URI") == "" {
		t.Skip("skipping: AGENTHOUND_PG_URI not set")
	}
}

func TestEnsureAdminUser_CreatesDefault(t *testing.T) {
	skipIfNoPG(t)
	ctx := context.Background()

	pool, err := appdb.NewPool(os.Getenv("AGENTHOUND_PG_URI"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := appdb.RunMigrations(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// Clean slate
	_, _ = pool.Exec(ctx, "DELETE FROM users")

	userStore := appdb.NewUserStore(pool)
	if err := EnsureAdminUser(ctx, userStore, "testpass"); err != nil {
		t.Fatalf("EnsureAdminUser: %v", err)
	}

	user, err := userStore.GetByUsername(ctx, "admin")
	if err != nil {
		t.Fatalf("GetByUsername: %v", err)
	}
	if user.Username != "admin" {
		t.Errorf("username: got %q, want admin", user.Username)
	}
	if user.Role != RoleAdmin {
		t.Errorf("role: got %q, want %q", user.Role, RoleAdmin)
	}
	if err := CheckPassword(user.PasswordHash, "testpass"); err != nil {
		t.Errorf("password mismatch: %v", err)
	}

	// Cleanup
	_, _ = pool.Exec(ctx, "DELETE FROM users WHERE username = 'admin'")
}

func TestEnsureAdminUser_Idempotent(t *testing.T) {
	skipIfNoPG(t)
	ctx := context.Background()

	pool, err := appdb.NewPool(os.Getenv("AGENTHOUND_PG_URI"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := appdb.RunMigrations(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM users")

	userStore := appdb.NewUserStore(pool)

	if err := EnsureAdminUser(ctx, userStore, "testpass"); err != nil {
		t.Fatalf("first call: %v", err)
	}

	if err := EnsureAdminUser(ctx, userStore, "otherpass"); err != nil {
		t.Fatalf("second call: %v", err)
	}

	count, err := userStore.Count(ctx)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("user count: got %d, want 1", count)
	}

	// Cleanup
	_, _ = pool.Exec(ctx, "DELETE FROM users WHERE username = 'admin'")
}
