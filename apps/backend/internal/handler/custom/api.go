package custom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/go-chi/chi/v5"
)

// MountCustomHandlers mounts all custom API endpoints
func MountCustomHandlers(r chi.Router, client *ent.Client) {
	r.Get("/status", getStatusHandler())
	r.Post("/login", loginHandler())

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware) // Apply authentication middleware to this group

		r.Get("/profile", getProfileHandler())
		r.Get("/users/count", getTotalUserCountHandler(client))

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireRoleMiddleware("admin"))
			r.Get("/admin/users/list", getAdminUserListHandler())
		})
	})
}

// getStatusHandler godoc
// @Summary Get API Status
// @Description Returns the current status of the API.
// @Tags Utilities
// @Produce json
// @Success 200 {object} custom.StatusResponse "API status is OK"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /status [get]
func getStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := StatusResponse{
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(resp)
	}
}

// loginHandler godoc
// @Summary User Login
// @Description Authenticates a user and returns an access token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body custom.LoginRequest true "User credentials"
// @Success 200 {object} custom.LoginResponse "Login successful"
// @Failure 400 {object} auth.BackendError "Invalid input"
// @Failure 401 {object} auth.BackendError "Invalid credentials"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /login [post]
func loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement actual login logic here
		// For example, parse LoginRequest from body, validate credentials, generate JWT
		w.Header().Set("Content-Type", "application/json")
		resp := LoginResponse{Token: "supersecrettoken-dummy-jwt"} // Dummy
		json.NewEncoder(w).Encode(resp)
	}
}

// getProfileHandler godoc
// @Summary Get User Profile
// @Description Returns the profile of the authenticated user.
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} auth.AuthenticatedUser "User profile"
// @Failure 401 {object} auth.BackendError "Unauthorized"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /profile [get]
func getProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok || user == nil {
			http.Error(w, "User not authenticated or found in context", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user) // Assuming AuthenticatedUser is JSON serializable
	}
}

// getTotalUserCountHandler godoc
// @Summary Get Total User Count
// @Description Returns the total number of users in the system.
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} custom.TotalUsersResponse "Total user count"
// @Failure 401 {object} auth.BackendError "Unauthorized"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /users/count [get]
func getTotalUserCountHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := client.User.Query().Count(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get user count: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := TotalUsersResponse{TotalUsers: count}
		json.NewEncoder(w).Encode(resp)
	}
}

// getAdminUserListHandler godoc
// @Summary Admin-Only User List
// @Description Returns a list of users (admin access only).
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "Admin-only user list!"
// @Failure 401 {object} auth.BackendError "Unauthorized"
// @Failure 403 {object} auth.BackendError "Forbidden"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /admin/users/list [get]
func getAdminUserListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Dummy admin list
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Admin-only user list!", "users": [{"id":1,"name":"Admin User"}]}`))
	}
}
