package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	// jwt "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user" // Key to store user info in context

// AuthenticatedUser swagger:model AuthenticatedUser
// AuthenticatedUser represents the user's information after authentication
// used by Swagger for API documentation
type AuthenticatedUser struct {
	ID    int      `json:"id" example:"1"`
	Email string   `json:"email" example:"admin@example.com"`
	Roles []string `json:"roles" example:"[\"admin\", \"user\"]"`
	// Add other relevant user data
}

// BackendError swagger:model BackendError
// BackendError is a standardized error response structure, referenced by Swagger
// used by Swagger for API documentation
type BackendError struct {
	Code    int         `json:"code" example:"400"`
	Status  string      `json:"status" example:"Bad Request"`
	Errors  interface{} `json:"errors,omitempty"`                                       // Can be map[string]string or []string or null
	Message string      `json:"message,omitempty" example:"Detailed error description"` // Custom message
}

// Dummy user lookup (replace with actual DB lookup or JWT validation)
func lookupUserByToken(token string) (*AuthenticatedUser, error) {
	// In a real app:
	// 1. Validate JWT token (check signature, expiry, claims)
	// 2. Extract user ID/email from claims
	// 3. Look up user details in your database (using Ent client)
	// For now, a dummy:
	if token == "supersecrettoken" {
		return &AuthenticatedUser{ID: 1, Email: "admin@example.com", Roles: []string{"admin", "user"}}, nil
	}
	if token == "usertoken" {
		return &AuthenticatedUser{ID: 2, Email: "test@example.com", Roles: []string{"user"}}, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// AuthMiddleware extracts and validates a JWT token from the Authorization header.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Authorization header required"})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Authorization header must be in 'Bearer <token>' format"})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := parts[1]

		user, err := lookupUserByToken(token) // Replace with actual JWT validation/DB lookup
		if err != nil {
			json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Invalid or expired token"})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Store the authenticated user in the request context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext retrieves the AuthenticatedUser from the request context
func GetUserFromContext(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(userContextKey).(*AuthenticatedUser)
	return user, ok
}

// RequireRoleMiddleware checks if the authenticated user has any of the required roles.
func RequireRoleMiddleware(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok || user == nil {
				json.NewEncoder(w).Encode(BackendError{Code: http.StatusUnauthorized, Status: "Unauthorized", Message: "Unauthorized: User not found in context"})
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range user.Roles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				json.NewEncoder(w).Encode(BackendError{Code: http.StatusForbidden, Status: "Forbidden", Message: "Forbidden: Insufficient permissions"})
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
