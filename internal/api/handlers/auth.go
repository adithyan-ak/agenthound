package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/adithyan-ak/agenthound/internal/appdb"
	"github.com/adithyan-ak/agenthound/internal/audit"
	"github.com/adithyan-ak/agenthound/internal/auth"
	"github.com/adithyan-ak/agenthound/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userStore  *appdb.UserStore
	tokenStore *appdb.TokenStore
	jwtSecret  string
	audit      *audit.Logger
}

func NewAuthHandler(userStore *appdb.UserStore, tokenStore *appdb.TokenStore, jwtSecret string, auditLog *audit.Logger) *AuthHandler {
	return &AuthHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		jwtSecret:  jwtSecret,
		audit:      auditLog,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      userInfo  `json:"user"`
}

type userInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		WriteValidationError(w, "username and password are required")
		return
	}

	user, err := h.userStore.GetByUsername(r.Context(), req.Username)
	if err != nil {
		h.auditLog(r, "auth.login_failure", map[string]any{"username": req.Username, "reason": "unknown_user"})
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		h.auditLog(r, "auth.login_failure", map[string]any{"username": req.Username, "reason": "bad_password"})
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	token, expiresAt, err := auth.GenerateJWT(user.ID, user.Username, user.Role, h.jwtSecret, 24*time.Hour)
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}

	_ = h.userStore.UpdateLastLogin(r.Context(), user.ID)
	h.auditLog(r, "auth.login_success", map[string]any{"user_id": user.ID, "username": user.Username})

	WriteJSON(w, http.StatusOK, loginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      userInfo{ID: user.ID, Username: user.Username, Role: user.Role},
	})
}

type createTokenRequest struct {
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type createTokenResponse struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

func (h *AuthHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "invalid request body")
		return
	}
	if req.Name == "" {
		WriteValidationError(w, "token name is required")
		return
	}

	raw, hash, err := auth.GenerateAPIToken()
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}

	apiToken := &model.APIToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: hash,
		Name:      req.Name,
		ExpiresAt: req.ExpiresAt,
	}

	if err := h.tokenStore.Create(r.Context(), apiToken); err != nil {
		WriteInternalError(w, r, err)
		return
	}

	h.auditLog(r, "auth.token_created", map[string]any{"token_id": apiToken.ID, "token_name": apiToken.Name})

	WriteJSON(w, http.StatusCreated, createTokenResponse{
		Token: raw,
		ID:    apiToken.ID,
		Name:  apiToken.Name,
	})
}

func (h *AuthHandler) HandleListTokens(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	tokens, err := h.tokenStore.ListByUser(r.Context(), user.ID)
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) HandleDeleteToken(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.tokenStore.Delete(r.Context(), id); err != nil {
		WriteInternalError(w, r, err)
		return
	}
	h.auditLog(r, "auth.token_deleted", map[string]any{"token_id": id})
	w.WriteHeader(http.StatusNoContent)
}

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *AuthHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		WriteValidationError(w, "username and password are required")
		return
	}

	validRoles := map[string]bool{"admin": true, "analyst": true, "viewer": true}
	if !validRoles[req.Role] {
		WriteValidationError(w, "role must be admin, analyst, or viewer")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: hash,
		Role:         req.Role,
	}

	if err := h.userStore.Create(r.Context(), user); err != nil {
		WriteError(w, http.StatusConflict, "CONFLICT", "username already exists")
		return
	}

	h.auditLog(r, "auth.user_created", map[string]any{"user_id": user.ID, "username": user.Username, "role": user.Role})

	WriteJSON(w, http.StatusCreated, userInfo{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (h *AuthHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.List(r.Context(), 100, 0)
	if err != nil {
		WriteInternalError(w, r, err)
		return
	}

	result := make([]userInfo, len(users))
	for i, u := range users {
		result[i] = userInfo{ID: u.ID, Username: u.Username, Role: u.Role}
	}
	WriteJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) auditLog(r *http.Request, action string, details map[string]any) {
	if h.audit == nil {
		return
	}
	if err := h.audit.Log(r.Context(), action, details); err != nil {
		slog.Warn("audit log failed", "action", action, "error", err)
	}
}
