package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/handler/apikey"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

// Common auth errors.
var (
	errMissingAuthHeader  = errors.New("authorization header required")
	errInvalidAuthFormat  = errors.New("authorization header must be in '<type> <token>' format")
	errUnsupportedAuth    = errors.New("unsupported authorization type")
	errInvalidCredentials = errors.New("invalid or expired token/key")
	errInactiveUser       = errors.New("user account is inactive")
	errInvalidBasicFormat = errors.New("invalid basic auth format")
	errEmailMismatch      = errors.New("email does not match API key owner")
)

// AuthMiddleware extracts and validates a JWT token or an API key from the Authorization header.
func AuthMiddleware(authSvc coreauth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authenticate(r, authSvc)
			if err != nil {
				httputil.WriteError(w, r, http.StatusUnauthorized, err.Error(), nil)
				return
			}

			if !user.User.IsActive {
				httputil.WriteError(w, r, http.StatusUnauthorized, errInactiveUser.Error(), nil)
				return
			}

			ctx := context.WithValue(r.Context(), httputil.ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// authenticate parses the Authorization header and validates the credentials.
func authenticate(r *http.Request, authSvc coreauth.Service) (*entities.UserAccount, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errMissingAuthHeader
	}

	authType, token, err := parseAuthHeader(authHeader)
	if err != nil {
		return nil, err
	}

	return validateCredentials(r.Context(), authSvc, authType, token)
}

// parseAuthHeader splits the Authorization header into type and token.
func parseAuthHeader(header string) (authType, token string, err error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", "", errInvalidAuthFormat
	}
	return strings.ToLower(parts[0]), parts[1], nil
}

// validateCredentials dispatches to the appropriate validation strategy based on auth type.
func validateCredentials(ctx context.Context, authSvc coreauth.Service, authType, token string) (*entities.UserAccount, error) {
	switch authType {
	case "bearer":
		return validateBearer(ctx, authSvc, token)
	case "apikey":
		return validateAPIKey(ctx, authSvc, token)
	case "basic":
		return validateBasicAuth(ctx, authSvc, token)
	default:
		return nil, errUnsupportedAuth
	}
}

// validateBearer validates a Bearer token (JWT or API key with prefix).
func validateBearer(ctx context.Context, authSvc coreauth.Service, token string) (*entities.UserAccount, error) {
	// Try JWT first
	user, err := authSvc.ValidateToken(token)
	if err == nil {
		return user, nil
	}

	// Fall back to API key if token has the expected prefix
	if strings.HasPrefix(token, apikey.APIKeyPrefix) {
		return authSvc.ValidateAPIKey(ctx, token)
	}

	return nil, errInvalidCredentials
}

// validateAPIKey validates an API key directly.
func validateAPIKey(ctx context.Context, authSvc coreauth.Service, token string) (*entities.UserAccount, error) {
	user, err := authSvc.ValidateAPIKey(ctx, token)
	if err != nil {
		return nil, errInvalidCredentials
	}
	return user, nil
}

// validateBasicAuth handles Basic Auth where username is email and password is API key.
func validateBasicAuth(ctx context.Context, authSvc coreauth.Service, encodedCredentials string) (*entities.UserAccount, error) {
	email, apiKey, err := decodeBasicCredentials(encodedCredentials)
	if err != nil {
		return nil, err
	}

	user, err := authSvc.ValidateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, errInvalidCredentials
	}

	if !strings.EqualFold(user.Person.Email, email) {
		return nil, errEmailMismatch
	}

	return user, nil
}

// decodeBasicCredentials decodes base64 credentials and extracts email and password.
func decodeBasicCredentials(encoded string) (email, password string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", errInvalidBasicFormat
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", errInvalidBasicFormat
	}

	return parts[0], parts[1], nil
}

// RequireSysAdminMiddleware checks if the authenticated user is a system administrator.
func RequireSysAdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := httputil.GetUserFromContext(r.Context())
			if !ok || user == nil {
				httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized: user not found in context", nil)
				return
			}

			if !user.User.IsSysAdmin {
				httputil.WriteError(w, r, http.StatusForbidden, "forbidden: insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
