// internal/approvals/policy_default.go
package approvals

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/approvals"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

// DefaultPolicy is the concrete mapping per event type.
type DefaultPolicy struct {
	byType map[string]approvals.EventApprovalPolicy
}

func NewDefaultPolicy() *DefaultPolicy {
	p := &DefaultPolicy{
		byType: map[string]approvals.EventApprovalPolicy{},
	}

	p.byType[events.ProjectRoleAssignedType] = approvals.EventApprovalPolicy{
		InitiatorProjectRoles: []string{"admin"},
		Approvers: []approvals.ApproverSpec{
			{Kind: approvals.ApproverProjectRole, RoleKey: "admin"},
			{Kind: approvals.ApproverOrgRole, RoleKey: "admin", BubbleStrategy: approvals.BubbleStrategyFirstAncestor},
		},
		Resolution: "any_one",
		Quorum:     nil,
	}

	p.byType[events.ProjectRoleUnassignedType] = approvals.EventApprovalPolicy{
		InitiatorProjectRoles: []string{"admin"},
		Approvers: []approvals.ApproverSpec{
			{Kind: approvals.ApproverProjectRole, RoleKey: "admin"},
			{Kind: approvals.ApproverOrgRole, RoleKey: "admin", BubbleStrategy: approvals.BubbleStrategyFirstAncestor},
		},
		Resolution: "any_one",
		Quorum:     nil,
	}

	p.byType[events.ProductAddedType] = approvals.EventApprovalPolicy{
		InitiatorProjectRoles: []string{"admin"},
		Approvers: []approvals.ApproverSpec{
			{Kind: approvals.ApproverProjectRole, RoleKey: "admin"},
			{Kind: approvals.ApproverOrgRole, RoleKey: "admin", BubbleStrategy: approvals.BubbleStrategyFirstAncestor},
		},
		Resolution: "any_one",
	}
	p.byType[events.ProductRemovedType] = approvals.EventApprovalPolicy{
		InitiatorProjectRoles: []string{"admin"},
		Approvers: []approvals.ApproverSpec{
			{Kind: approvals.ApproverProjectRole, RoleKey: "admin"},
			{Kind: approvals.ApproverOrgRole, RoleKey: "admin", BubbleStrategy: approvals.BubbleStrategyFirstAncestor},
		},
		Resolution: "any_one",
	}

	return p
}

func (p *DefaultPolicy) ForEventType(eventType string) (approvals.EventApprovalPolicy, bool) {
	ap, ok := p.byType[eventType]
	return ap, ok
}
