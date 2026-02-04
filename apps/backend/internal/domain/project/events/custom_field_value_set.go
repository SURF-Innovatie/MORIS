package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	projdomain "github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const CustomFieldValueSetType = "project.custom_field_value_set"

type CustomFieldValueSet struct {
	Base
	DefinitionID string `json:"definition_id"`
	Value        string `json:"value"`
}

type CustomFieldValueSetInput struct {
	DefinitionID string `json:"definition_id"`
	Value        string `json:"value"`
}

func (CustomFieldValueSet) isEvent() {}

func (e CustomFieldValueSet) Action() string {
	return "custom_field_value_set"
}

func (CustomFieldValueSet) Type() string { return CustomFieldValueSetType }

func (e CustomFieldValueSet) String() string {
	return fmt.Sprintf("custom field %s set to %s", e.DefinitionID, e.Value)
}

func (e CustomFieldValueSet) Apply(p *projdomain.Project) {
	if p.CustomFields == nil {
		p.CustomFields = make(map[string]interface{})
	}
	p.CustomFields[e.DefinitionID] = e.Value
}

// Decider
func DecideCustomFieldValueSet(ctx context.Context, projectID, userID uuid.UUID, state *projdomain.Project, cmd CustomFieldValueSetInput, status Status) (Event, error) {
	if state == nil {
		return nil, errors.New("project does not exist")
	}

	if cmd.DefinitionID == "" {
		return nil, errors.New("definition_id is required")
	}

	base := NewBase(projectID, userID, status)
	base.FriendlyNameStr = CustomFieldValueSetMeta.FriendlyName

	return &CustomFieldValueSet{
		Base:         base,
		DefinitionID: cmd.DefinitionID,
		Value:        cmd.Value,
	}, nil
}

// Meta
var CustomFieldValueSetMeta = EventMeta{
	Type:         CustomFieldValueSetType,
	FriendlyName: "Set Custom Field Value",
	CheckAllowed: func(ctx context.Context, e Event, cli *ent.Client) bool { return true },
}
