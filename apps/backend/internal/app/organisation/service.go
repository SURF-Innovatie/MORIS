package organisation

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service interface {
	CreateRoot(ctx context.Context, name string, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error)
	CreateChild(ctx context.Context, parentID uuid.UUID, name string, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error)

	Get(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error)
	Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error)

	ListRoots(ctx context.Context) ([]entities.OrganisationNode, error)
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error)
	UpdateMemberCustomFields(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, values map[string]interface{}) error
}

type service struct {
	repo       Repository
	personRepo PersonRepository
}

func NewService(repo Repository, personRepo PersonRepository) Service {
	return &service{repo: repo, personRepo: personRepo}
}

func (s *service) CreateRoot(ctx context.Context, name string, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error) {
	var out *entities.OrganisationNode
	err := s.repo.WithTx(ctx, func(ctx context.Context, tx Repository) error {
		row, err := tx.CreateNode(ctx, name, nil, rorID, description, avatarURL)
		if err != nil {
			return err
		}
		if err := tx.InsertClosure(ctx, row.ID, row.ID, 0); err != nil {
			return err
		}
		out = row
		return nil
	})
	return out, err
}

func (s *service) CreateChild(ctx context.Context, parentID uuid.UUID, name string, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error) {
	var out *entities.OrganisationNode
	err := s.repo.WithTx(ctx, func(ctx context.Context, tx Repository) error {
		// ensure parent exists
		if _, err := tx.GetNode(ctx, parentID); err != nil {
			return err
		}

		row, err := tx.CreateNode(ctx, name, &parentID, rorID, description, avatarURL)
		if err != nil {
			return err
		}

		// self closure
		if err := tx.InsertClosure(ctx, row.ID, row.ID, 0); err != nil {
			return err
		}

		// ancestor closures from parent's ancestors
		ancRows, err := tx.ListClosuresByDescendant(ctx, parentID)
		if err != nil {
			return err
		}
		bulk := lo.Map(ancRows, func(a entities.OrganisationNodeClosure, _ int) entities.OrganisationNodeClosure {
			return entities.OrganisationNodeClosure{
				AncestorID:   a.AncestorID,
				DescendantID: row.ID,
				Depth:        a.Depth + 1,
			}
		})
		if len(bulk) > 0 {
			if err := tx.CreateClosuresBulk(ctx, bulk); err != nil {
				return err
			}
		}

		out = row
		return nil
	})
	return out, err
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error) {
	return s.repo.GetNode(ctx, id)
}

func (s *service) ListRoots(ctx context.Context) ([]entities.OrganisationNode, error) {
	return s.repo.ListRoots(ctx)
}

func (s *service) ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error) {
	return s.repo.ListChildren(ctx, parentID)
}

func (s *service) UpdateMemberCustomFields(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, values map[string]interface{}) error {
	p, err := s.personRepo.Get(ctx, personID)
	if err != nil {
		return err
	}

	if p.OrgCustomFields == nil {
		p.OrgCustomFields = make(map[string]interface{})
	}

	p.OrgCustomFields[orgID.String()] = values

	_, err = s.personRepo.Update(ctx, personID, *p)
	return err
}

func (s *service) Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error) {
	var out *entities.OrganisationNode
	err := s.repo.WithTx(ctx, func(ctx context.Context, tx Repository) error {
		cur, err := tx.GetNode(ctx, id)
		if err != nil {
			return err
		}
		if parentID != nil && *parentID == id {
			return fmt.Errorf("node cannot be its own parent")
		}

		row, err := tx.UpdateNode(ctx, id, name, parentID, rorID, description, avatarURL)
		if err != nil {
			return err
		}

		// If parent didn't change -> done
		oldParentID := cur.ParentID
		newParentID := parentID

		same :=
			(oldParentID == nil && newParentID == nil) ||
				(oldParentID != nil && newParentID != nil && *oldParentID == *newParentID)
		if same {
			out = row
			return nil
		}

		// subtree: closures where ancestor=id
		subRows, err := tx.ListClosuresByAncestor(ctx, id)
		if err != nil {
			return err
		}
		subDepth := lo.Associate(subRows, func(r entities.OrganisationNodeClosure) (uuid.UUID, int) {
			return r.DescendantID, r.Depth
		})
		subIDs := lo.Keys(subDepth)

		// cycle check: cannot move under own subtree
		if newParentID != nil {
			if _, ok := subDepth[*newParentID]; ok {
				return fmt.Errorf("cannot move node under its own subtree")
			}
		}

		// old external ancestors: closures where descendant=id AND ancestor NOT IN subtree
		oldAncRows, err := tx.ListClosuresByDescendant(ctx, id)
		if err != nil {
			return err
		}
		subSet := lo.SliceToMap(subIDs, func(sid uuid.UUID) (uuid.UUID, struct{}) {
			return sid, struct{}{}
		})
		oldAncIDs := lo.FilterMap(oldAncRows, func(a entities.OrganisationNodeClosure, _ int) (uuid.UUID, bool) {
			if _, inSub := subSet[a.AncestorID]; !inSub {
				return a.AncestorID, true
			}
			return uuid.Nil, false
		})

		// delete: (old external ancestors) -> (subtree descendants)
		if len(oldAncIDs) > 0 && len(subIDs) > 0 {
			if err := tx.DeleteClosures(ctx, oldAncIDs, subIDs); err != nil {
				return err
			}
		}

		// new root? done
		if newParentID == nil {
			out = row
			return nil
		}

		// new external ancestors: closures where descendant=newParent
		newAncRows, err := tx.ListClosuresByDescendant(ctx, *newParentID)
		if err != nil {
			return err
		}

		bulk := lo.FlatMap(newAncRows, func(na entities.OrganisationNodeClosure, _ int) []entities.OrganisationNodeClosure {
			return lo.Map(subIDs, func(descID uuid.UUID, _ int) entities.OrganisationNodeClosure {
				return entities.OrganisationNodeClosure{
					AncestorID:   na.AncestorID,
					DescendantID: descID,
					Depth:        na.Depth + 1 + subDepth[descID],
				}
			})
		})
		if len(bulk) > 0 {
			if err := tx.CreateClosuresBulk(ctx, bulk); err != nil {
				return err
			}
		}

		out = row
		return nil
	})
	return out, err
}
