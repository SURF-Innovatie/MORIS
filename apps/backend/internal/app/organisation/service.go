package organisation

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	rbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/tx"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/readmodels"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service interface {
	CreateRoot(ctx context.Context, name string, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error)
	CreateChild(ctx context.Context, parentID uuid.UUID, name string, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error)

	Get(ctx context.Context, id uuid.UUID) (*organisation.OrganisationNode, error)
	Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error)

	ListRoots(ctx context.Context) ([]organisation.OrganisationNode, error)
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]organisation.OrganisationNode, error)
	ListAll(ctx context.Context) ([]organisation.OrganisationNode, error)
	Search(ctx context.Context, query string) ([]organisation.OrganisationNode, error)
	SearchForProjectCreation(ctx context.Context, query string, actorID uuid.UUID) ([]organisation.OrganisationNode, error)
	UpdateMemberCustomFields(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, values map[string]interface{}) error
}

type service struct {
	repo      repository
	personSvc person.Service
	rbac      rbacsvc.Service
	tx        tx.Manager
}

func NewService(repo repository, personRepo person.Service, rbac rbacsvc.Service, tx tx.Manager) Service {
	return &service{repo: repo, personSvc: personRepo, rbac: rbac, tx: tx}
}

func (s *service) CreateRoot(ctx context.Context, name string, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error) {
	var out *organisation.OrganisationNode
	slug := generateSlug(name)
	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		row, err := s.repo.CreateNode(ctx, name, nil, rorID, description, avatarURL, slug)
		if err != nil {
			return err
		}
		if err := s.repo.InsertClosure(ctx, row.ID, row.ID, 0); err != nil {
			return err
		}
		out = row
		return nil
	})
	return out, err
}

func (s *service) CreateChild(ctx context.Context, parentID uuid.UUID, name string, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error) {
	var out *organisation.OrganisationNode
	slug := generateSlug(name)
	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		// ensure parent exists
		if _, err := s.repo.GetNode(ctx, parentID); err != nil {
			return err
		}

		row, err := s.repo.CreateNode(ctx, name, &parentID, rorID, description, avatarURL, slug)
		if err != nil {
			return err
		}

		// self closure
		if err := s.repo.InsertClosure(ctx, row.ID, row.ID, 0); err != nil {
			return err
		}

		// ancestor closures from parent's ancestors
		ancRows, err := s.repo.ListClosuresByDescendant(ctx, parentID)
		if err != nil {
			return err
		}
		bulk := lo.Map(ancRows, func(a readmodels.OrganisationNodeClosure, _ int) readmodels.OrganisationNodeClosure {
			return readmodels.OrganisationNodeClosure{
				AncestorID:   a.AncestorID,
				DescendantID: row.ID,
				Depth:        a.Depth + 1,
			}
		})
		if len(bulk) > 0 {
			if err := s.repo.CreateClosuresBulk(ctx, bulk); err != nil {
				return err
			}
		}

		out = row
		return nil
	})
	return out, err
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*organisation.OrganisationNode, error) {
	return s.repo.GetNode(ctx, id)
}

func (s *service) ListRoots(ctx context.Context) ([]organisation.OrganisationNode, error) {
	return s.repo.ListRoots(ctx)
}

func (s *service) ListChildren(ctx context.Context, parentID uuid.UUID) ([]organisation.OrganisationNode, error) {
	return s.repo.ListChildren(ctx, parentID)
}

func (s *service) ListAll(ctx context.Context) ([]organisation.OrganisationNode, error) {
	return s.repo.ListAll(ctx)
}

func (s *service) Search(ctx context.Context, query string) ([]organisation.OrganisationNode, error) {
	return s.repo.Search(ctx, query, 20)
}

func (s *service) SearchForProjectCreation(ctx context.Context, query string, actorID uuid.UUID) ([]organisation.OrganisationNode, error) {
	nodes, err := s.repo.Search(ctx, query, 100)
	if err != nil {
		return nil, err
	}

	// Filter nodes where the user has PermissionCreateProject
	// HasPermission checks for direct or inherited permission (via ancestry)
	filtered := make([]organisation.OrganisationNode, 0, len(nodes))
	for _, node := range nodes {
		ok, err := s.rbac.HasPermission(ctx, actorID, node.ID, rbac.PermissionCreateProject)
		if err == nil && ok {
			filtered = append(filtered, node)
		}
	}
	return filtered, nil
}

func (s *service) UpdateMemberCustomFields(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, values map[string]interface{}) error {
	p, err := s.personSvc.Get(ctx, personID)
	if err != nil {
		return err
	}

	if p.OrgCustomFields == nil {
		p.OrgCustomFields = make(map[string]interface{})
	}

	p.OrgCustomFields[orgID.String()] = values

	_, err = s.personSvc.Update(ctx, personID, *p)
	return err
}

func (s *service) Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error) {
	var out *organisation.OrganisationNode
	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		cur, err := s.repo.GetNode(ctx, id)
		if err != nil {
			return err
		}
		if parentID != nil && *parentID == id {
			return fmt.Errorf("node cannot be its own parent")
		}

		row, err := s.repo.UpdateNode(ctx, id, name, parentID, rorID, description, avatarURL)
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
		subRows, err := s.repo.ListClosuresByAncestor(ctx, id)
		if err != nil {
			return err
		}
		subDepth := lo.Associate(subRows, func(r readmodels.OrganisationNodeClosure) (uuid.UUID, int) {
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
		oldAncRows, err := s.repo.ListClosuresByDescendant(ctx, id)
		if err != nil {
			return err
		}
		subSet := lo.SliceToMap(subIDs, func(sid uuid.UUID) (uuid.UUID, struct{}) {
			return sid, struct{}{}
		})
		oldAncIDs := lo.FilterMap(oldAncRows, func(a readmodels.OrganisationNodeClosure, _ int) (uuid.UUID, bool) {
			if _, inSub := subSet[a.AncestorID]; !inSub {
				return a.AncestorID, true
			}
			return uuid.Nil, false
		})

		// delete: (old external ancestors) -> (subtree descendants)
		if len(oldAncIDs) > 0 && len(subIDs) > 0 {
			if err := s.repo.DeleteClosures(ctx, oldAncIDs, subIDs); err != nil {
				return err
			}
		}

		// new root? done
		if newParentID == nil {
			out = row
			return nil
		}

		// new external ancestors: closures where descendant=newParent
		newAncRows, err := s.repo.ListClosuresByDescendant(ctx, *newParentID)
		if err != nil {
			return err
		}

		bulk := lo.FlatMap(newAncRows, func(na readmodels.OrganisationNodeClosure, _ int) []readmodels.OrganisationNodeClosure {
			return lo.Map(subIDs, func(descID uuid.UUID, _ int) readmodels.OrganisationNodeClosure {
				return readmodels.OrganisationNodeClosure{
					AncestorID:   na.AncestorID,
					DescendantID: descID,
					Depth:        na.Depth + 1 + subDepth[descID],
				}
			})
		})
		if len(bulk) > 0 {
			if err := s.repo.CreateClosuresBulk(ctx, bulk); err != nil {
				return err
			}
		}

		out = row
		return nil
	})
	return out, err
}

func generateSlug(name string) string {
	return strings.ToLower(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(name, "-"))
}
