package auth

import "net/http"

const (
	RoleAdmin   = "admin"
	RoleAnalyst = "analyst"
	RoleViewer  = "viewer"
)

var roleHierarchy = map[string]int{
	RoleViewer:  1,
	RoleAnalyst: 2,
	RoleAdmin:   3,
}

func RequireRole(minRole string) func(http.Handler) http.Handler {
	minLevel := roleHierarchy[minRole]
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
				return
			}

			userLevel := roleHierarchy[user.Role]
			if userLevel < minLevel {
				writeAuthError(w, http.StatusForbidden, "FORBIDDEN", "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
