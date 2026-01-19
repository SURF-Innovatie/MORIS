package middleware

import (
	"context"
	"net/http"
	"strings"

	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

// AuthMiddleware extracts and validates a JWT token or an API key from the Authorization header.
func AuthMiddleware(authSvc coreauth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.WriteError(w, r, http.StatusUnauthorized, "Authorization header required", nil)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				httputil.WriteError(w, r, http.StatusUnauthorized, "Authorization header must be in '<type> <token>' format", nil)
				return
			}

			authType := strings.ToLower(parts[0])
			token := parts[1]

			var user *entities.UserAccount
			var err error

			if authType == "bearer" {
				// Try JWT first
				user, err = authSvc.ValidateToken(token)
				if err != nil {
					// If JWT fails, check if it's an API key (starts with moris_)
					// TODO: refine api key implementation etc
					if strings.HasPrefix(token, "moris_") {
						user, err = authSvc.ValidateAPIKey(r.Context(), token)
					}
				}
			} else if authType == "apikey" {
				user, err = authSvc.ValidateAPIKey(r.Context(), token)
			} else {
				httputil.WriteError(w, r, http.StatusUnauthorized, "Unsupported authorization type", nil)
				return
			}

			if err != nil {
				httputil.WriteError(w, r, http.StatusUnauthorized, "Invalid or expired token/key", nil)
				return
			}

			// Check if user is active
			if !user.User.IsActive {
				httputil.WriteError(w, r, http.StatusUnauthorized, "User account is inactive", nil)
				return
			}

			// Store the authenticated user in the request context
			ctx := context.WithValue(r.Context(), httputil.ContextKeyUser, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSysAdminMiddleware checks if the authenticated user is a system administrator.
func RequireSysAdminMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := httputil.GetUserFromContext(r.Context())
			if !ok || user == nil {
				httputil.WriteError(w, r, http.StatusUnauthorized, "Unauthorized: User not found in context", nil)
				return
			}

			if !user.User.IsSysAdmin {
				httputil.WriteError(w, r, http.StatusForbidden, "Forbidden: Insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
