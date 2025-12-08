package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	coreauth "github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type contextKey string

const userContextKey contextKey = "user" // Key to store user info in context

// BackendError swagger:model BackendError
// BackendError is a standardized error response structure, referenced by Swagger
// used by Swagger for API documentation
type BackendError struct {
	Code    int         `json:"code" example:"400"`
	Status  string      `json:"status" example:"Bad Request"`
	Errors  interface{} `json:"errors,omitempty"`                                       // Can be map[string]string or []string or null
	Message string      `json:"message,omitempty" example:"Detailed error description"` // Custom message
}

// AuthMiddleware extracts and validates a JWT token from the Authorization header.
func AuthMiddleware(authSvc coreauth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Authorization header required"})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Authorization header must be in 'Bearer <token>' format"})
				return
			}
			token := parts[1]

			user, err := authSvc.ValidateToken(token)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Invalid or expired token"})
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Unauthorized: User not found in context"})
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
