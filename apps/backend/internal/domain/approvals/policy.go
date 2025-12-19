package approvals

type ApproverKind string

const (
	ApproverProjectRole ApproverKind = "project_role"
	ApproverOrgRole     ApproverKind = "org_role"
)

type BubbleStrategy int

const (
	BubbleStrategyFirstAncestor = iota
	BubbleStrategyAllAncestors
)

var BubbleStrategyMap = map[BubbleStrategy]string{
	BubbleStrategyFirstAncestor: "first_ancestor",
	BubbleStrategyAllAncestors:  "all_ancestors",
}

func (bs BubbleStrategy) String() string {
	return BubbleStrategyMap[bs]
}

// ApproverSpec defines a role that can approve an event
type ApproverSpec struct {
	// Kind indicates whether the approver is a project role or an org role
	Kind ApproverKind
	// RoleKey is the key of the role that serves as approver
	RoleKey string
	// BubbleStrategy indicates how to bubble up this approver through project hierarchy
	// Only relevant when the Approver is a ApproverOrgRole
	BubbleStrategy BubbleStrategy
}

type EventApprovalPolicy struct {
	InitiatorProjectRoles []string // who may create the request at all
	Approvers             []ApproverSpec
	Resolution            string // "any_one" | "all" | "quorum"
	Quorum                *int
}

type Policy interface {
	ForEventType(eventType string) (EventApprovalPolicy, bool)
}
