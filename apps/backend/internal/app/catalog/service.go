package catalog

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/catalog"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	ErrNotFound = errors.New("catalog not found")
)

type CatalogDetails struct {
	Catalog       ent.Catalog                                 `json:"catalog"`
	Projects      []*queries.ProjectDetails                   `json:"projects"`
	People        map[uuid.UUID]identity.Person               `json:"people"`
	Products      map[uuid.UUID]product.Product               `json:"products"`
	Organisations map[uuid.UUID]organisation.OrganisationNode `json:"organisations"`
}

type CreateRequest struct {
	Name            string      `json:"name"`
	Title           string      `json:"title"`
	Description     *string     `json:"description"`
	RichDescription *string     `json:"rich_description"`
	ProjectIDs      []uuid.UUID `json:"project_ids"`
	LogoURL         *string     `json:"logo_url"`
	PrimaryColor    *string     `json:"primary_color"`
	SecondaryColor  *string     `json:"secondary_color"`
	AccentColor     *string     `json:"accent_color"`
	Favicon         *string     `json:"favicon"`
	FontFamily      *string     `json:"font_family"`
}

type UpdateRequest struct {
	Name            *string      `json:"name"`
	Title           *string      `json:"title"`
	Description     *string      `json:"description"`
	RichDescription *string      `json:"rich_description"`
	ProjectIDs      *[]uuid.UUID `json:"project_ids"`
	LogoURL         *string      `json:"logo_url"`
	PrimaryColor    *string      `json:"primary_color"`
	SecondaryColor  *string      `json:"secondary_color"`
	AccentColor     *string      `json:"accent_color"`
	Favicon         *string      `json:"favicon"`
	FontFamily      *string      `json:"font_family"`
}

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*ent.Catalog, error)
	Get(ctx context.Context, id uuid.UUID) (*ent.Catalog, error)
	GetDetails(ctx context.Context, id uuid.UUID) (*CatalogDetails, error)
	List(ctx context.Context) ([]*ent.Catalog, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*ent.Catalog, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	client     *ent.Client
	projectSvc queries.Service
}

func NewService(client *ent.Client, projectSvc queries.Service) Service {
	return &service{
		client:     client,
		projectSvc: projectSvc,
	}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*ent.Catalog, error) {
	builder := s.client.Catalog.Create().
		SetName(req.Name).
		SetTitle(req.Title).
		SetProjectIds(req.ProjectIDs)

	if req.Description != nil {
		builder.SetDescription(*req.Description)
	}
	if req.RichDescription != nil {
		builder.SetRichDescription(*req.RichDescription)
	}
	if req.LogoURL != nil {
		builder.SetLogoURL(*req.LogoURL)
	}
	if req.PrimaryColor != nil {
		builder.SetPrimaryColor(*req.PrimaryColor)
	}
	if req.SecondaryColor != nil {
		builder.SetSecondaryColor(*req.SecondaryColor)
	}
	if req.AccentColor != nil {
		builder.SetAccentColor(*req.AccentColor)
	}
	if req.Favicon != nil {
		builder.SetFavicon(*req.Favicon)
	}
	if req.FontFamily != nil {
		builder.SetFontFamily(*req.FontFamily)
	}

	return builder.Save(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*ent.Catalog, error) {
	return s.client.Catalog.Query().Where(catalog.ID(id)).Only(ctx)
}

func (s *service) List(ctx context.Context) ([]*ent.Catalog, error) {
	return s.client.Catalog.Query().All(ctx)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*ent.Catalog, error) {
	builder := s.client.Catalog.UpdateOneID(id)

	if req.Name != nil {
		builder.SetName(*req.Name)
	}
	if req.Title != nil {
		builder.SetTitle(*req.Title)
	}
	if req.Description != nil {
		builder.SetDescription(*req.Description)
	}
	if req.RichDescription != nil {
		builder.SetRichDescription(*req.RichDescription)
	}
	if req.ProjectIDs != nil {
		builder.SetProjectIds(*req.ProjectIDs)
	}
	if req.LogoURL != nil {
		builder.SetLogoURL(*req.LogoURL)
	}
	if req.PrimaryColor != nil {
		builder.SetPrimaryColor(*req.PrimaryColor)
	}
	if req.SecondaryColor != nil {
		builder.SetSecondaryColor(*req.SecondaryColor)
	}
	if req.AccentColor != nil {
		builder.SetAccentColor(*req.AccentColor)
	}
	if req.Favicon != nil {
		builder.SetFavicon(*req.Favicon)
	}
	if req.FontFamily != nil {
		builder.SetFontFamily(*req.FontFamily)
	}

	return builder.Save(ctx)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.client.Catalog.DeleteOneID(id).Exec(ctx)
}

func (s *service) GetDetails(ctx context.Context, id uuid.UUID) (*CatalogDetails, error) {
	cat, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Fetch all projects details
	log.Debug().Msgf("GetDetails: found %d project IDs for catalog %s", len(cat.ProjectIds), id)
	projects := make([]*queries.ProjectDetails, 0, len(cat.ProjectIds))
	for _, pid := range cat.ProjectIds {
		// Use a context that allows access
		p, err := s.projectSvc.GetProject(ctx, pid)
		if err == nil {
			projects = append(projects, p)
		} else {
			log.Error().Err(err).Msgf("GetDetails: failed to get project %s", pid)
		}
	}

	// Aggregate entities
	peopleIDs := make([]uuid.UUID, 0)
	productIDs := make([]uuid.UUID, 0)

	for _, p := range projects {
		for _, m := range p.Members {
			// Person uses ID (all caps)
			peopleIDs = append(peopleIDs, m.Person.ID)
		}
		for _, prod := range p.Products {
			// Product uses Id (PascalCase)
			productIDs = append(productIDs, prod.Id)
		}
	}

	uniquePeopleIDs := lo.Uniq(peopleIDs)
	uniqueProductIDs := lo.Uniq(productIDs)

	people, _ := s.projectSvc.GetPeopleByIDs(ctx, uniquePeopleIDs)
	products, _ := s.projectSvc.GetProductsByIDs(ctx, uniqueProductIDs)

	// Temporarily: we have OrgNode in ProjectDetails, so we can aggregate from there.
	orgs := make(map[uuid.UUID]organisation.OrganisationNode)
	for _, p := range projects {
		orgs[p.OwningOrgNode.ID] = p.OwningOrgNode
	}

	return &CatalogDetails{
		Catalog:       *cat,
		Projects:      projects,
		People:        people,
		Products:      products,
		Organisations: orgs,
	}, nil
}
