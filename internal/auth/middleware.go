package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/adithyan-ak/agenthound/internal/appdb"
	"github.com/adithyan-ak/agenthound/internal/model"
)

type Middleware struct {
	jwtSecret  string
	tokenStore *appdb.TokenStore
	userStore  *appdb.UserStore
}

func NewMiddleware(jwtSecret string, tokenStore *appdb.TokenStore, userStore *appdb.UserStore) *Middleware {
	return &Middleware{
		jwtSecret:  jwtSecret,
		tokenStore: tokenStore,
		userStore:  userStore,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token == "" {
			writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing or invalid authorization header")
			return
		}

		user, err := m.resolveUser(r, token)
		if err != nil {
			writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) resolveUser(r *http.Request, token string) (*model.User, error) {
	if strings.HasPrefix(token, tokenPrefix) {
		return m.resolveAPIToken(r, token)
	}
	return m.resolveJWT(token)
}

func (m *Middleware) resolveJWT(token string) (*model.User, error) {
	claims, err := ValidateJWT(token, m.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

func (m *Middleware) resolveAPIToken(r *http.Request, raw string) (*model.User, error) {
	hash := HashToken(raw)
	apiToken, err := m.tokenStore.GetByHash(r.Context(), hash)
	if err != nil {
		return nil, err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := m.tokenStore.UpdateLastUsed(ctx, apiToken.ID); err != nil {
			slog.Warn("failed to update token last_used", "token_id", apiToken.ID, "error", err)
		}
	}()

	user, err := m.userStore.GetByID(r.Context(), apiToken.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func writeAuthError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":{"code":"` + code + `","message":"` + message + `"}}`))
}
