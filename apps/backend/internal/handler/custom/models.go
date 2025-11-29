package custom

import (
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/google/uuid"
)

// RegisterRequest swagger:model RegisterRequest
// Represents the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secretpassword"`
}

// RegisterResponse swagger:model RegisterResponse
// Represents the response body for successful registration.
type RegisterResponse struct {
	ID       uuid.UUID `json:"id" example:"1"`
	PersonID uuid.UUID `json:"person_id" example:"2"`
}

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
	User  *userdto.Response
}

// StatusResponse swagger:model StatusResponse
// Represents the /status endpoint response payload.
type StatusResponse struct {
	Status    string `json:"status" example:"ok"`
	Timestamp string `json:"timestamp" example:"2025-11-12T10:00:00Z"`
}

// ORCIDAuthURLResponse swagger:model ORCIDAuthURLResponse
// Represents the response body for getting the ORCID auth URL.
type ORCIDAuthURLResponse struct {
	URL string `json:"url" example:"https://orcid.org/oauth/authorize?..."`
}

// LinkORCIDRequest swagger:model LinkORCIDRequest
// Represents the request body for linking an ORCID ID.
type LinkORCIDRequest struct {
	Code string `json:"code" example:"authentication_code_from_orcid"`
}
