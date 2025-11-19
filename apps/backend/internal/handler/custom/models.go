package custom

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
	ID    int    `json:"id" example:"1"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
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
	User  struct {
		ID    int      `json:"id" example:"1"`
		Email string   `json:"email" example:"user@example.com"`
		Roles []string `json:"roles" example:"user"`
	} `json:"user"`
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

// Note: BackendError and AuthenticatedUser types are defined in the
// `internal/auth` package, and annotated with `swagger:model`. We don't need
// wrapper types here anymore; referencing the canonical `auth` types in
// handler comments will prevent duplicate schema definitions in the OpenAPI spec.
