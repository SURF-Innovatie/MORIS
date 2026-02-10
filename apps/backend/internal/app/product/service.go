package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*product.Product, error)
	GetAll(ctx context.Context) ([]*product.Product, error)
	GetAllForUser(ctx context.Context, personID uuid.UUID) ([]*product.Product, error)
	Create(ctx context.Context, p product.Product) (*product.Product, error)
	Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByDOI(ctx context.Context, doi string) (*product.Product, error)
	CreateOrGetFromWork(ctx context.Context, work *dto.Work) (*product.Product, bool, error)
	GetOrCreateFromDOI(ctx context.Context, doi string) (*product.Product, bool, error)
}

type service struct {
	repo      Repository
	doiSvc    doi.Service
	personSvc person.Service
}

func NewService(repo Repository, doiSvc doi.Service, personSvc person.Service) Service {
	return &service{repo: repo, doiSvc: doiSvc, personSvc: personSvc}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]*product.Product, error) {
	return s.repo.List(ctx)
}

func (s *service) GetAllForUser(ctx context.Context, personID uuid.UUID) ([]*product.Product, error) {
	return s.repo.ListByAuthorPersonID(ctx, personID)
}

func (s *service) Create(ctx context.Context, p product.Product) (*product.Product, error) {
	return s.repo.Create(ctx, p)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error) {
	return s.repo.Update(ctx, id, p)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetByDOI(ctx context.Context, doi string) (*product.Product, error) {
	doi = strings.TrimSpace(doi)
	if doi == "" {
		return nil, nil
	}
	return s.repo.GetByDOI(ctx, doi)
}

func (s *service) CreateOrGetFromWork(ctx context.Context, w *dto.Work) (*product.Product, bool, error) {
	if w == nil {
		return nil, false, fmt.Errorf("work is required")
	}
	doiStr := strings.TrimSpace(w.DOI.String())
	if doiStr == "" {
		return nil, false, fmt.Errorf("work.doi is required")
	}

	existing, err := s.repo.GetByDOI(ctx, doiStr)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, false, nil
	}

	authorIDs, err := s.resolveAuthorPersons(ctx, w.Authors)
	if err != nil {
		return nil, false, err
	}

	p := product.Product{
		Name:            w.Title,
		Type:            w.Type,
		Language:        "en",
		DOI:             doiStr,
		AuthorPersonIDs: authorIDs,
	}

	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, false, err
	}
	return created, true, nil
}

func (s *service) GetOrCreateFromDOI(ctx context.Context, doiStr string) (*product.Product, bool, error) {
	doiStr = strings.TrimSpace(doiStr)
	if doiStr == "" {
		return nil, false, fmt.Errorf("doi is required")
	}

	existing, err := s.repo.GetByDOI(ctx, doiStr)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, false, nil
	}

	if s.doiSvc == nil {
		return nil, false, fmt.Errorf("doi service not configured")
	}

	work, err := s.doiSvc.Resolve(ctx, doiStr)
	if err != nil {
		return nil, false, err
	}

	return s.CreateOrGetFromWork(ctx, work)
}

func (s *service) resolveAuthorPersons(ctx context.Context, authors []dto.WorkAuthor) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(authors))

	for _, a := range authors {
		orcid := strings.TrimSpace(a.ORCID)
		name := strings.TrimSpace(a.Name)
		given := strings.TrimSpace(a.Given)
		family := strings.TrimSpace(a.Family)

		var p *identity.Person
		var err error

		if orcid != "" {
			p, err = s.personSvc.GetByORCID(ctx, orcid)
			if err != nil {
				return nil, err
			}
		}

		// If ORCID absent or not found -> always create a new person (per your requirement)
		if p == nil {
			if name == "" {
				if given != "" && family != "" {
					name = given + " " + family
				} else if family != "" {
					name = family
				} else {
					name = "Unknown author"
				}
			}

			email := makeSyntheticAuthorEmail(orcid, name)

			created, err := s.personSvc.Create(ctx, identity.Person{
				ORCiD:      ptrIfNonEmpty(orcid), // depending on your identity.Person type
				Name:       name,
				GivenName:  ptrIfNonEmpty(given),
				FamilyName: ptrIfNonEmpty(family),
				Email:      email,
			})
			if err != nil {
				return nil, err
			}
			p = created
		}

		ids = append(ids, p.ID) // adjust field name to your identity.Person
	}

	// optional: dedupe ids
	uniq := make([]uuid.UUID, 0, len(ids))
	seen := map[uuid.UUID]struct{}{}
	for _, id := range ids {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	return uniq, nil
}

func makeSyntheticAuthorEmail(orcid string, fallbackName string) string {
	// Stable + unique enough for your system (adjust domain)
	key := orcid
	if key == "" {
		key = uuid.NewString()
	}
	return fmt.Sprintf("author+%s@moris.invalid", key)
}

func ptrIfNonEmpty(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
