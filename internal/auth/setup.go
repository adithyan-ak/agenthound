package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adithyan-ak/agenthound/internal/appdb"
	"github.com/adithyan-ak/agenthound/internal/model"
	"github.com/google/uuid"
)

func EnsureAdminUser(ctx context.Context, userStore *appdb.UserStore, defaultPassword string) error {
	count, err := userStore.Count(ctx)
	if err != nil {
		return fmt.Errorf("check user count: %w", err)
	}
	if count > 0 {
		return nil
	}

	hash, err := HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     "admin",
		PasswordHash: hash,
		Role:         RoleAdmin,
	}

	if err := userStore.Create(ctx, user); err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	slog.Warn("created default admin user", "username", "admin")
	if defaultPassword == "agenthound" {
		slog.Warn("using default admin password — change it or set AGENTHOUND_ADMIN_PASSWORD")
	}

	return nil
}
