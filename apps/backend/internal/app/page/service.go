package page

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/types"
	"github.com/google/uuid"
)

var (
	ErrPageNotFound      = errors.New("page not found")
	ErrUserPageExists    = errors.New("user already has a page")
	ErrProjectPageExists = errors.New("project already has a page")
)

type Service interface {
	CreatePage(ctx context.Context, title, slug, pType string, content []types.Section, projectID, userID *uuid.UUID) (*ent.Page, error)
	GetPage(ctx context.Context, slug string) (*ent.Page, error)
	UpdatePage(ctx context.Context, id uuid.UUID, title *string, content []types.Section, isPublished *bool) (*ent.Page, error)
	ListProjectPages(ctx context.Context, projectID uuid.UUID) ([]*ent.Page, error)
	ListUserPages(ctx context.Context, userID uuid.UUID) ([]*ent.Page, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePage(ctx context.Context, title, slug, pType string, content []types.Section, projectID, userID *uuid.UUID) (*ent.Page, error) {
	// Enforce one page per user
	if userID != nil {
		existingPages, err := s.repo.ListByUser(ctx, *userID)
		if err != nil {
			return nil, err
		}
		if len(existingPages) > 0 {
			return nil, ErrUserPageExists
		}
	}

	// Enforce one page per project
	if projectID != nil {
		existingPages, err := s.repo.ListByProject(ctx, *projectID)
		if err != nil {
			return nil, err
		}
		if len(existingPages) > 0 {
			return nil, ErrProjectPageExists
		}
	}

	return s.repo.Create(ctx, title, slug, pType, content, projectID, userID)
}

func (s *service) GetPage(ctx context.Context, slug string) (*ent.Page, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *service) UpdatePage(ctx context.Context, id uuid.UUID, title *string, content []types.Section, isPublished *bool) (*ent.Page, error) {
	return s.repo.Update(ctx, id, title, content, isPublished)
}

func (s *service) ListProjectPages(ctx context.Context, projectID uuid.UUID) ([]*ent.Page, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *service) ListUserPages(ctx context.Context, userID uuid.UUID) ([]*ent.Page, error) {
	return s.repo.ListByUser(ctx, userID)
}
