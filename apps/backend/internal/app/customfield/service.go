package customfield

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/customfield"
	"github.com/google/uuid"
)

var (
	ErrNotFoundInOrg = errors.New("custom_field_definition_not_found_in_organization")
)

type Service interface {
	Create(ctx context.Context, orgID uuid.UUID, name string, fieldType customfield.Type, category customfield.Category, description, validationRegex, exampleValue *string, required bool) (*customfield.Definition, error)
	Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
	ListAvailableForNode(ctx context.Context, orgID uuid.UUID, category *customfield.Category) ([]customfield.Definition, error)
	ListAvailableForProject(ctx context.Context, projectID uuid.UUID) ([]*customfield.Definition, error)
}

type CreateDefinitionInput struct {
	OrgID           uuid.UUID
	Name            string
	Type            customfield.Type
	Category        customfield.Category
	Description     *string
	ValidationRegex *string
	ExampleValue    *string
	Required        bool
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(
	ctx context.Context,
	orgID uuid.UUID,
	name string,
	fieldType customfield.Type,
	category customfield.Category,
	description, validationRegex, exampleValue *string,
	required bool,
) (*customfield.Definition, error) {
	return s.repo.Create(ctx, CreateDefinitionInput{
		OrgID:           orgID,
		Name:            name,
		Type:            fieldType,
		Category:        category,
		Description:     description,
		ValidationRegex: validationRegex,
		ExampleValue:    exampleValue,
		Required:        required,
	})
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	exists, err := s.repo.ExistsInOrg(ctx, id, orgID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w", ErrNotFoundInOrg)
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) ListAvailableForNode(ctx context.Context, orgID uuid.UUID, category *customfield.Category) ([]customfield.Definition, error) {
	return s.repo.ListAvailableForNode(ctx, orgID, category)
}

func (s *service) ListAvailableForProject(ctx context.Context, projectID uuid.UUID) ([]*customfield.Definition, error) {
	return nil, fmt.Errorf("not implemented, use ListAvailableForNode with project's orgID")
}
