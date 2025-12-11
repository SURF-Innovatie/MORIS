package entities

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Organisation struct {
	Id   uuid.UUID
	Name string
}

func (o *Organisation) FromEnt(row *ent.Organisation) *Organisation {
	return &Organisation{
		Id:   row.ID,
		Name: row.Name,
	}
}
