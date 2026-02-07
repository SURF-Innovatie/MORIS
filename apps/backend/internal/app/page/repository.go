package page

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/page"
	"github.com/SURF-Innovatie/MORIS/internal/types"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, title, slug, pType string, content []types.Section, projectID, userID *uuid.UUID) (*ent.Page, error)
	Get(ctx context.Context, id uuid.UUID) (*ent.Page, error)
	GetBySlug(ctx context.Context, slug string) (*ent.Page, error)
	Update(ctx context.Context, id uuid.UUID, title *string, content []types.Section, isPublished *bool) (*ent.Page, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*ent.Page, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*ent.Page, error)
}

type repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	return &repository{client: client}
}

func (r *repository) Create(ctx context.Context, title, slug, pType string, content []types.Section, projectID, userID *uuid.UUID) (*ent.Page, error) {
	query := r.client.Page.Create().
		SetTitle(title).
		SetSlug(slug).
		SetType(page.Type(pType)).
		SetContent(content)

	if projectID != nil {
		query.SetProjectID(*projectID)
	}
	if userID != nil {
		query.SetUserID(*userID)
	}

	return query.Save(ctx)
}

func (r *repository) Get(ctx context.Context, id uuid.UUID) (*ent.Page, error) {
	return r.client.Page.Get(ctx, id)
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (*ent.Page, error) {
	return r.client.Page.Query().
		Where(page.Slug(slug)).
		Only(ctx)
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, title *string, content []types.Section, isPublished *bool) (*ent.Page, error) {
	update := r.client.Page.UpdateOneID(id)
	if title != nil {
		update.SetTitle(*title)
	}
	if content != nil {
		update.SetContent(content)
	}
	if isPublished != nil {
		update.SetIsPublished(*isPublished)
	}
	return update.Save(ctx)
}

func (r *repository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*ent.Page, error) {
	return r.client.Page.Query().
		Where(page.ProjectID(projectID)).
		All(ctx)
}

func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*ent.Page, error) {
	return r.client.Page.Query().
		Where(page.UserID(userID)).
		All(ctx)
}
