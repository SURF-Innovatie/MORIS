package organisationdto

import "github.com/google/uuid"

type Response struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
