package eventpolicy

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/policy"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type EntRepo struct {
	client *ent.Client
}

// NewEntRepository creates a new ent-based event policy repository
func NewEntRepository(client *ent.Client) *EntRepo {
	return &EntRepo{client: client}
}

func (r *EntRepo) Create(ctx context.Context, eventPolicy policy.EventPolicy) (*policy.EventPolicy, error) {
	create := r.client.EventPolicy.Create().
		SetName(eventPolicy.Name).
		SetEventTypes(eventPolicy.EventTypes).
		SetActionType(eventpolicy.ActionType(eventPolicy.ActionType)).
		SetEnabled(eventPolicy.Enabled)

	if eventPolicy.Description != nil {
		create.SetDescription(*eventPolicy.Description)
	}
	if len(eventPolicy.Conditions) > 0 {
		create.SetConditions(eventPolicy.ConditionsToMap())
	}
	if eventPolicy.MessageTemplate != nil {
		create.SetMessageTemplate(*eventPolicy.MessageTemplate)
	}
	if len(eventPolicy.RecipientUserIDs) > 0 {
		create.SetRecipientUserIds(eventPolicy.RecipientUserIDs)
	}
	if len(eventPolicy.RecipientProjectRoleIDs) > 0 {
		create.SetRecipientProjectRoleIds(eventPolicy.RecipientProjectRoleIDs)
	}
	if len(eventPolicy.RecipientOrgRoleIDs) > 0 {
		create.SetRecipientOrgRoleIds(eventPolicy.RecipientOrgRoleIDs)
	}
	if len(eventPolicy.RecipientDynamic) > 0 {
		create.SetRecipientDynamic(eventPolicy.RecipientDynamic)
	}
	if eventPolicy.OrgNodeID != nil {
		create.SetOrgNodeID(*eventPolicy.OrgNodeID)
	}
	if eventPolicy.ProjectID != nil {
		create.SetProjectID(*eventPolicy.ProjectID)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return new(policy.EventPolicy).FromEnt(row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, eventPolicy policy.EventPolicy) (*policy.EventPolicy, error) {
	update := r.client.EventPolicy.UpdateOneID(id).
		SetName(eventPolicy.Name).
		SetEventTypes(eventPolicy.EventTypes).
		SetActionType(eventpolicy.ActionType(eventPolicy.ActionType)).
		SetEnabled(eventPolicy.Enabled)

	if eventPolicy.Description != nil {
		update.SetDescription(*eventPolicy.Description)
	} else {
		update.ClearDescription()
	}

	if len(eventPolicy.Conditions) > 0 {
		update.SetConditions(eventPolicy.ConditionsToMap())
	} else {
		update.ClearConditions()
	}

	if eventPolicy.MessageTemplate != nil {
		update.SetMessageTemplate(*eventPolicy.MessageTemplate)
	} else {
		update.ClearMessageTemplate()
	}

	if len(eventPolicy.RecipientUserIDs) > 0 {
		update.SetRecipientUserIds(eventPolicy.RecipientUserIDs)
	} else {
		update.ClearRecipientUserIds()
	}

	if len(eventPolicy.RecipientProjectRoleIDs) > 0 {
		update.SetRecipientProjectRoleIds(eventPolicy.RecipientProjectRoleIDs)
	} else {
		update.ClearRecipientProjectRoleIds()
	}

	if len(eventPolicy.RecipientOrgRoleIDs) > 0 {
		update.SetRecipientOrgRoleIds(eventPolicy.RecipientOrgRoleIDs)
	} else {
		update.ClearRecipientOrgRoleIds()
	}

	if len(eventPolicy.RecipientDynamic) > 0 {
		update.SetRecipientDynamic(eventPolicy.RecipientDynamic)
	} else {
		update.ClearRecipientDynamic()
	}

	row, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return new(policy.EventPolicy).FromEnt(row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.EventPolicy.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) GetByID(ctx context.Context, id uuid.UUID) (*policy.EventPolicy, error) {
	row, err := r.client.EventPolicy.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return new(policy.EventPolicy).FromEnt(row), nil
}

func (r *EntRepo) ListForOrgNode(ctx context.Context, orgNodeID uuid.UUID, ancestorNodeIDs []uuid.UUID) ([]policy.EventPolicy, error) {
	// Build query for org node ID or any of its ancestors
	allIDs := append([]uuid.UUID{orgNodeID}, ancestorNodeIDs...)

	rows, err := r.client.EventPolicy.Query().
		Where(eventpolicy.OrgNodeIDIn(allIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(row *ent.EventPolicy, _ int) policy.EventPolicy {
		return *new(policy.EventPolicy).FromEnt(row)
	}), nil
}

func (r *EntRepo) ListForProject(ctx context.Context, projectID uuid.UUID) ([]policy.EventPolicy, error) {
	rows, err := r.client.EventPolicy.Query().
		Where(eventpolicy.ProjectIDEQ(projectID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(row *ent.EventPolicy, _ int) policy.EventPolicy {
		return *new(policy.EventPolicy).FromEnt(row)
	}), nil
}
