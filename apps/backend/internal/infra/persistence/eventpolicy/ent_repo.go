package eventpolicy

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type entRepo struct {
	client *ent.Client
}

// NewEntRepository creates a new ent-based event policy repository
func NewEntRepository(client *ent.Client) *entRepo {
	return &entRepo{client: client}
}

func (r *entRepo) Create(ctx context.Context, policy entities.EventPolicy) (*entities.EventPolicy, error) {
	create := r.client.EventPolicy.Create().
		SetName(policy.Name).
		SetEventTypes(policy.EventTypes).
		SetActionType(eventpolicy.ActionType(policy.ActionType)).
		SetEnabled(policy.Enabled)

	if policy.Description != nil {
		create.SetDescription(*policy.Description)
	}
	if len(policy.Conditions) > 0 {
		create.SetConditions(policy.ConditionsToMap())
	}
	if policy.MessageTemplate != nil {
		create.SetMessageTemplate(*policy.MessageTemplate)
	}
	if len(policy.RecipientUserIDs) > 0 {
		create.SetRecipientUserIds(policy.RecipientUserIDs)
	}
	if len(policy.RecipientProjectRoleIDs) > 0 {
		create.SetRecipientProjectRoleIds(policy.RecipientProjectRoleIDs)
	}
	if len(policy.RecipientOrgRoleIDs) > 0 {
		create.SetRecipientOrgRoleIds(policy.RecipientOrgRoleIDs)
	}
	if len(policy.RecipientDynamic) > 0 {
		create.SetRecipientDynamic(policy.RecipientDynamic)
	}
	if policy.OrgNodeID != nil {
		create.SetOrgNodeID(*policy.OrgNodeID)
	}
	if policy.ProjectID != nil {
		create.SetProjectID(*policy.ProjectID)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return new(entities.EventPolicy).FromEnt(row), nil
}

func (r *entRepo) Update(ctx context.Context, id uuid.UUID, policy entities.EventPolicy) (*entities.EventPolicy, error) {
	update := r.client.EventPolicy.UpdateOneID(id).
		SetName(policy.Name).
		SetEventTypes(policy.EventTypes).
		SetActionType(eventpolicy.ActionType(policy.ActionType)).
		SetEnabled(policy.Enabled)

	if policy.Description != nil {
		update.SetDescription(*policy.Description)
	} else {
		update.ClearDescription()
	}

	if len(policy.Conditions) > 0 {
		update.SetConditions(policy.ConditionsToMap())
	} else {
		update.ClearConditions()
	}

	if policy.MessageTemplate != nil {
		update.SetMessageTemplate(*policy.MessageTemplate)
	} else {
		update.ClearMessageTemplate()
	}

	if len(policy.RecipientUserIDs) > 0 {
		update.SetRecipientUserIds(policy.RecipientUserIDs)
	} else {
		update.ClearRecipientUserIds()
	}

	if len(policy.RecipientProjectRoleIDs) > 0 {
		update.SetRecipientProjectRoleIds(policy.RecipientProjectRoleIDs)
	} else {
		update.ClearRecipientProjectRoleIds()
	}

	if len(policy.RecipientOrgRoleIDs) > 0 {
		update.SetRecipientOrgRoleIds(policy.RecipientOrgRoleIDs)
	} else {
		update.ClearRecipientOrgRoleIds()
	}

	if len(policy.RecipientDynamic) > 0 {
		update.SetRecipientDynamic(policy.RecipientDynamic)
	} else {
		update.ClearRecipientDynamic()
	}

	row, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return new(entities.EventPolicy).FromEnt(row), nil
}

func (r *entRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.EventPolicy.DeleteOneID(id).Exec(ctx)
}

func (r *entRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.EventPolicy, error) {
	row, err := r.client.EventPolicy.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return new(entities.EventPolicy).FromEnt(row), nil
}

func (r *entRepo) ListForOrgNode(ctx context.Context, orgNodeID uuid.UUID, ancestorNodeIDs []uuid.UUID) ([]entities.EventPolicy, error) {
	// Build query for org node ID or any of its ancestors
	allIDs := append([]uuid.UUID{orgNodeID}, ancestorNodeIDs...)

	rows, err := r.client.EventPolicy.Query().
		Where(eventpolicy.OrgNodeIDIn(allIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(row *ent.EventPolicy, _ int) entities.EventPolicy {
		return *new(entities.EventPolicy).FromEnt(row)
	}), nil
}

func (r *entRepo) ListForProject(ctx context.Context, projectID uuid.UUID) ([]entities.EventPolicy, error) {
	rows, err := r.client.EventPolicy.Query().
		Where(eventpolicy.ProjectIDEQ(projectID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(row *ent.EventPolicy, _ int) entities.EventPolicy {
		return *new(entities.EventPolicy).FromEnt(row)
	}), nil
}
