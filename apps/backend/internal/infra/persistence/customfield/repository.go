package customfield

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/customfielddefinition"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type repo struct {
	cli *ent.Client
}

func NewRepository(cli *ent.Client) customfield.Repository {
	return &repo{cli: cli}
}

func (r *repo) Create(ctx context.Context, in customfield.CreateDefinitionInput) (*entities.CustomFieldDefinition, error) {
	creator := r.cli.CustomFieldDefinition.Create().
		SetOrganisationNodeID(in.OrgID).
		SetName(in.Name).
		SetType(customfielddefinition.Type(in.Type)).
		SetCategory(customfielddefinition.Category(in.Category)).
		SetRequired(in.Required)

	if in.Description != nil {
		creator.SetDescription(*in.Description)
	}
	if in.ValidationRegex != nil {
		creator.SetValidationRegex(*in.ValidationRegex)
	}
	if in.ExampleValue != nil {
		creator.SetExampleValue(*in.ExampleValue)
	}

	row, err := creator.Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.CustomFieldDefinition](row), nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.CustomFieldDefinition.DeleteOneID(id).Exec(ctx)
}

func (r *repo) ExistsInOrg(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (bool, error) {
	return r.cli.CustomFieldDefinition.Query().
		Where(
			customfielddefinition.ID(id),
			customfielddefinition.OrganisationNodeID(orgID),
		).
		Exist(ctx)
}

func (r *repo) ListAvailableForNode(ctx context.Context, orgID uuid.UUID, category *entities.CustomFieldCategory) ([]entities.CustomFieldDefinition, error) {
	// Use closure table to find ancestors (including self if closure contains it)
	ancestors, err := r.cli.OrganisationNodeClosure.Query().
		Where(organisationnodeclosure.DescendantIDEQ(orgID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting ancestors for org %s: %w", orgID, err)
	}

	validOrgIDs := make([]uuid.UUID, 0, len(ancestors))
	for _, a := range ancestors {
		validOrgIDs = append(validOrgIDs, a.AncestorID)
	}

	q := r.cli.CustomFieldDefinition.Query().
		Where(customfielddefinition.OrganisationNodeIDIn(validOrgIDs...))

	if category != nil {
		q.Where(customfielddefinition.CategoryEQ(customfielddefinition.Category(*category)))
	}

	rows, err := q.Order(ent.Asc(customfielddefinition.FieldName)).All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.CustomFieldDefinition](rows), nil
}
