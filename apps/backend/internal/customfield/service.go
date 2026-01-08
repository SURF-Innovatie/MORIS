package customfield

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/customfielddefinition"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, orgID uuid.UUID, name string, fieldType customfielddefinition.Type, category customfielddefinition.Category, description, validationRegex, exampleValue *string, required bool) (*ent.CustomFieldDefinition, error)
	Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
	ListAvailableForNode(ctx context.Context, orgID uuid.UUID, category *customfielddefinition.Category) ([]*ent.CustomFieldDefinition, error)
	ListAvailableForProject(ctx context.Context, projectID uuid.UUID) ([]*ent.CustomFieldDefinition, error)
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, orgID uuid.UUID, name string, fieldType customfielddefinition.Type, category customfielddefinition.Category, description, validationRegex, exampleValue *string, required bool) (*ent.CustomFieldDefinition, error) {
	creator := s.cli.CustomFieldDefinition.Create().
		SetOrganisationNodeID(orgID).
		SetName(name).
		SetType(fieldType).
		SetCategory(category).
		SetRequired(required)

	if description != nil {
		creator.SetDescription(*description)
	}
	if validationRegex != nil {
		creator.SetValidationRegex(*validationRegex)
	}
	if exampleValue != nil {
		creator.SetExampleValue(*exampleValue)
	}

	return creator.Save(ctx)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	// Ensure the definition belongs to the org
	exists, err := s.cli.CustomFieldDefinition.Query().
		Where(
			customfielddefinition.ID(id),
			customfielddefinition.OrganisationNodeID(orgID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("custom field definition not found in organization")
	}

	return s.cli.CustomFieldDefinition.DeleteOneID(id).Exec(ctx)
}

func (s *service) ListAvailableForNode(ctx context.Context, orgID uuid.UUID, category *customfielddefinition.Category) ([]*ent.CustomFieldDefinition, error) {
	// Use Closure table to find ancestors
	ancestorIDs, err := s.cli.OrganisationNodeClosure.Query().
		Where(organisationnodeclosure.DescendantIDEQ(orgID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting ancestors for org %s: %w", orgID, err)
	}

	validOrgIDs := make([]uuid.UUID, 0, len(ancestorIDs))
	for _, a := range ancestorIDs {
		validOrgIDs = append(validOrgIDs, a.AncestorID)
	}

	// Query definitions
	q := s.cli.CustomFieldDefinition.Query().
		Where(customfielddefinition.OrganisationNodeIDIn(validOrgIDs...))

	if category != nil {
		q.Where(customfielddefinition.CategoryEQ(*category))
	}

	return q.Order(ent.Asc(customfielddefinition.FieldName)).
		All(ctx)
}

func (s *service) ListAvailableForProject(ctx context.Context, projectID uuid.UUID) ([]*ent.CustomFieldDefinition, error) {
	return nil, fmt.Errorf("not implemented, use ListAvailableForNode with project's orgID")
}


// Helper to get user from context if needed in future
func currentUser(ctx context.Context, cli *ent.Client) (*ent.User, error) {
	authUser, ok := httputil.GetUserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no authenticated user in context")
	}
	return cli.User.Get(ctx, authUser.User.ID)
}
