package portfolio

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entportfolio "github.com/SURF-Innovatie/MORIS/ent/portfolio"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/portfolio"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) GetByPersonID(ctx context.Context, personID uuid.UUID) (*portfolio.Portfolio, error) {
	row, err := r.cli.Portfolio.Query().
		Where(entportfolio.PersonIDEQ(personID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return transform.ToEntityPtr[portfolio.Portfolio](row), nil
}

func (r *EntRepo) Upsert(ctx context.Context, port portfolio.Portfolio) (*portfolio.Portfolio, error) {
	existing, err := r.cli.Portfolio.Query().
		Where(entportfolio.PersonIDEQ(port.PersonID)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if existing == nil {
		row, err := r.cli.Portfolio.Create().
			SetPersonID(port.PersonID).
			SetNillableHeadline(port.Headline).
			SetNillableSummary(port.Summary).
			SetNillableWebsite(port.Website).
			SetShowEmail(port.ShowEmail).
			SetShowOrcid(port.ShowOrcid).
			SetPinnedProjectIds(port.PinnedProjectIDs).
			SetPinnedProductIds(port.PinnedProductIDs).
			SetRecentProjectIds(port.RecentProjectIDs).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return transform.ToEntityPtr[portfolio.Portfolio](row), nil
	}

	row, err := r.cli.Portfolio.UpdateOne(existing).
		SetNillableHeadline(port.Headline).
		SetNillableSummary(port.Summary).
		SetNillableWebsite(port.Website).
		SetShowEmail(port.ShowEmail).
		SetShowOrcid(port.ShowOrcid).
		SetPinnedProjectIds(port.PinnedProjectIDs).
		SetPinnedProductIds(port.PinnedProductIDs).
		SetRecentProjectIds(port.RecentProjectIDs).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[portfolio.Portfolio](row), nil
}
