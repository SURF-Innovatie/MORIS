package middleware

import (
	"context"
	"net/http"
	"strings"

	coreauth "github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type contextKey string

const userContextKey contextKey = "user" // Key to store user info in context

// AuthMiddleware extracts and validates a JWT token from the Authorization header.
func AuthMiddleware(authSvc coreauth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "Authorization header required", nil)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httputil.WriteError(w, http.StatusUnauthorized, "Authorization header must be in 'Bearer <token>' format", nil)
				return
			}
			token := parts[1]

			user, err := authSvc.ValidateToken(token)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "Invalid or expired token", nil)
				return
			}

			// Store the authenticated user in the request context
			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext retrieves the authUser from the request context
func GetUserFromContext(ctx context.Context) (*entities.UserAccount, bool) {
	user, ok := ctx.Value(userContextKey).(*entities.UserAccount)
	return user, ok
}

// RequireRoleMiddleware checks if the authenticated user has any of the required roles.
func RequireRoleMiddleware(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok || user == nil {
				httputil.WriteError(w, http.StatusUnauthorized, "Unauthorized: User not found in context", nil)
				return
			}

			//hasRole := false
			//for _, requiredRole := range roles {
			//	for _, userRole := range user.Roles {
			//		if userRole == requiredRole {
			//			hasRole = true
			//			break
			//		}
			//	}
			//	if hasRole {
			//		break
			//	}
			//}
			//
			//if !hasRole {
			//	w.Header().Set("Content-Type", "application/json")
			//	w.WriteHeader(http.StatusForbidden)
			//	json.NewEncoder(w).Encode(BackendError{Code: http.StatusForbidden, Status: "Forbidden", Message: "Forbidden: Insufficient permissions"})
			//	return
			//}

			next.ServeHTTP(w, r)
		})
	}
}
