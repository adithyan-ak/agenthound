package auth

import (
	"context"

	"github.com/adithyan-ak/agenthound/internal/model"
)

type contextKey struct{}

func ContextWithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, contextKey{}, user)
}

func UserFromContext(ctx context.Context) *model.User {
	u, _ := ctx.Value(contextKey{}).(*model.User)
	return u
}
