package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"valid", "Bearer abc123", "abc123"},
		{"case insensitive", "bearer abc123", "abc123"},
		{"empty", "", ""},
		{"no bearer prefix", "Basic abc123", ""},
		{"bearer only", "Bearer", ""},
		{"spaces trimmed", "Bearer  abc123 ", "abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
			}
			got := extractBearerToken(r)
			if got != tt.want {
				t.Errorf("extractBearerToken() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAuthenticateNoHeader(t *testing.T) {
	m := &Middleware{jwtSecret: "secret"}
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestAuthenticateValidJWT(t *testing.T) {
	secret := "test-secret"
	token, _, err := GenerateJWT("user-1", "testuser", "admin", secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	m := &Middleware{jwtSecret: secret}
	var gotUser string
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r.Context())
		if u != nil {
			gotUser = u.Username
		}
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if gotUser != "testuser" {
		t.Errorf("user = %q, want testuser", gotUser)
	}
}

func TestAuthenticateExpiredJWT(t *testing.T) {
	secret := "test-secret"
	token, _, err := GenerateJWT("user-1", "testuser", "admin", secret, -time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	m := &Middleware{jwtSecret: secret}
	handler := m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestRequireRoleAllowed(t *testing.T) {
	tests := []struct {
		name    string
		minRole string
		user    string
		role    string
		want    int
	}{
		{"admin accessing admin route", RoleAdmin, "admin", RoleAdmin, 200},
		{"analyst accessing analyst route", RoleAnalyst, "analyst", RoleAnalyst, 200},
		{"admin accessing analyst route", RoleAnalyst, "admin", RoleAdmin, 200},
		{"admin accessing viewer route", RoleViewer, "admin", RoleAdmin, 200},
		{"viewer accessing viewer route", RoleViewer, "viewer", RoleViewer, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireRole(tt.minRole)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(ContextWithUser(r.Context(), &model.User{Username: tt.user, Role: tt.role}))
			handler.ServeHTTP(w, r)

			if w.Code != tt.want {
				t.Errorf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}

func TestRequireRoleDenied(t *testing.T) {
	tests := []struct {
		name    string
		minRole string
		role    string
	}{
		{"viewer accessing admin route", RoleAdmin, RoleViewer},
		{"viewer accessing analyst route", RoleAnalyst, RoleViewer},
		{"analyst accessing admin route", RoleAdmin, RoleAnalyst},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireRole(tt.minRole)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("handler should not be called")
			}))

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(ContextWithUser(r.Context(), &model.User{Role: tt.role}))
			handler.ServeHTTP(w, r)

			if w.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403", w.Code)
			}
		})
	}
}

func TestRequireRoleNoUser(t *testing.T) {
	handler := RequireRole(RoleViewer)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}
