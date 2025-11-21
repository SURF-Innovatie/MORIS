package organisationdto

import "github.com/google/uuid"

type CreateRequest struct {
	Name string `json:"name"`
}

type Response struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
