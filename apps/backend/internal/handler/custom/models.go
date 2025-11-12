package custom

import "github.com/SURF-Innovatie/MORIS/internal/auth"

// LoginRequest swagger:model LoginRequest
// Represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secretpassword"`
}

// LoginResponse swagger:model LoginResponse
// Represents the response body for successful login.
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// StatusResponse swagger:model StatusResponse
// Represents the /status endpoint response payload.
type StatusResponse struct {
	Status    string `json:"status" example:"ok"`
	Timestamp string `json:"timestamp" example:"2025-11-12T10:00:00Z"`
}

// TotalUsersResponse swagger:model TotalUsersResponse
// Represents the payload returned by /users/count.
type TotalUsersResponse struct {
	TotalUsers int `json:"total_users" example:"123"`
}

// BackendErrorDoc swagger:model BackendError
// Provides a clean schema name for backend error responses.
type BackendErrorDoc auth.BackendError

// AuthenticatedUserDoc swagger:model AuthenticatedUser
// Provides a clean schema name for authenticated user payloads.
type AuthenticatedUserDoc auth.AuthenticatedUser
