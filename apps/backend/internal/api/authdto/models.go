package authdto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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
	Token string                `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  *entities.UserAccount `json:"user"`
}

func FromEntity(token string, user *entities.UserAccount) LoginResponse {
	return LoginResponse{
		Token: token,
		User:  user,
	}
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
