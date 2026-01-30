package event_policy

import (
	"context"

	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/google/uuid"
)

// OrgClosureAdapter implements eventpolicy.OrgClosureProvider using the org repository
type OrgClosureAdapter struct {
	orgRbac organisationrbac.Service
}

// NewOrgClosureAdapter creates a new OrgClosureAdapter
func NewOrgClosureAdapter(orgRbac organisationrbac.Service) *OrgClosureAdapter {
	return &OrgClosureAdapter{orgRbac: orgRbac}
}

// GetAncestorIDs returns all ancestor org node IDs for a given node (excluding self)
func (a *OrgClosureAdapter) GetAncestorIDs(ctx context.Context, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	ids, err := a.orgRbac.AncestorIDs(ctx, orgNodeID)
	if err != nil {
		return nil, err
	}

	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if id != orgNodeID {
			out = append(out, id)
		}
	}
	return out, nil
}
