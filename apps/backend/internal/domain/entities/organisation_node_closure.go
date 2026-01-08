package entities

import "github.com/google/uuid"

type OrganisationNodeClosure struct {
	AncestorID   uuid.UUID
	DescendantID uuid.UUID
	Depth        int
}
