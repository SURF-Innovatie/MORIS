// Package hydrator provides centralized event enrichment logic.
// It hydrates events with related entities (Person, Product, ProjectRole, etc.)
package hydrator

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Loaders define the repository interfaces needed for hydration
type PersonLoader interface {
	PeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Person, error)
}

type ProductLoader interface {
	ProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]entities.Product, error)
}

type RoleLoader interface {
	ProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.ProjectRole, error)
}

type UserPersonResolver interface {
	GetPeopleByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]entities.Person, error)
}

// Hydrator enriches events with related entities
type Hydrator struct {
	persons  PersonLoader
	products ProductLoader
	roles    RoleLoader
	users    UserPersonResolver
}

// New creates a new Hydrator
func New(persons PersonLoader, products ProductLoader, roles RoleLoader, users UserPersonResolver) *Hydrator {
	return &Hydrator{
		persons:  persons,
		products: products,
		roles:    roles,
		users:    users,
	}
}

// HydrateOne hydrates a single event with related entities
func (h *Hydrator) HydrateOne(ctx context.Context, e events.Event) events.DetailedEvent {
	return h.HydrateMany(ctx, []events.Event{e})[0]
}

// HydrateMany hydrates multiple events with batch loading for efficiency
func (h *Hydrator) HydrateMany(ctx context.Context, evts []events.Event) []events.DetailedEvent {
	if len(evts) == 0 {
		return nil
	}

	// Collect all IDs to batch load
	var personIDs, roleIDs, productIDs, creatorUserIDs []uuid.UUID

	for _, e := range evts {
		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				personIDs = append(personIDs, *ids.PersonID)
			}
			if ids.ProjectRoleID != nil {
				roleIDs = append(roleIDs, *ids.ProjectRoleID)
			}
			if ids.ProductID != nil {
				productIDs = append(productIDs, *ids.ProductID)
			}
		}
		creatorUserIDs = append(creatorUserIDs, e.CreatedByID())
	}

	// Batch load all entities
	personMap := h.loadPersons(ctx, personIDs)
	roleMap := h.loadRoles(ctx, roleIDs)
	productMap := h.loadProducts(ctx, productIDs)
	creatorMap := h.loadCreators(ctx, creatorUserIDs)

	// Build detailed events
	return lo.Map(evts, func(e events.Event, _ int) events.DetailedEvent {
		de := events.DetailedEvent{Event: e}

		if r, ok := e.(events.HasRelatedIDs); ok {
			ids := r.RelatedIDs()
			if ids.PersonID != nil {
				if p, ok := personMap[*ids.PersonID]; ok {
					de.Person = &p
				}
			}
			if ids.ProjectRoleID != nil {
				if r, ok := roleMap[*ids.ProjectRoleID]; ok {
					de.ProjectRole = &r
				}
			}
			if ids.ProductID != nil {
				if p, ok := productMap[*ids.ProductID]; ok {
					de.Product = &p
				}
			}
		}

		if p, ok := creatorMap[e.CreatedByID()]; ok {
			de.Creator = &p
		}

		return de
	})
}

func (h *Hydrator) loadPersons(ctx context.Context, ids []uuid.UUID) map[uuid.UUID]entities.Person {
	if len(ids) == 0 {
		return nil
	}
	m, err := h.persons.PeopleByIDs(ctx, lo.Uniq(ids))
	if err != nil {
		return nil
	}
	return m
}

func (h *Hydrator) loadRoles(ctx context.Context, ids []uuid.UUID) map[uuid.UUID]entities.ProjectRole {
	if len(ids) == 0 {
		return nil
	}
	m, err := h.roles.ProjectRolesByIDs(ctx, lo.Uniq(ids))
	if err != nil {
		return nil
	}
	return m
}

func (h *Hydrator) loadProducts(ctx context.Context, ids []uuid.UUID) map[uuid.UUID]entities.Product {
	if len(ids) == 0 {
		return nil
	}
	products, err := h.products.ProductsByIDs(ctx, lo.Uniq(ids))
	if err != nil {
		return nil
	}
	return lo.SliceToMap(products, func(p entities.Product) (uuid.UUID, entities.Product) {
		return p.Id, p
	})
}

func (h *Hydrator) loadCreators(ctx context.Context, ids []uuid.UUID) map[uuid.UUID]entities.Person {
	if len(ids) == 0 {
		return nil
	}
	m, err := h.users.GetPeopleByUserIDs(ctx, lo.Uniq(ids))
	if err != nil {
		return nil
	}
	return m
}
