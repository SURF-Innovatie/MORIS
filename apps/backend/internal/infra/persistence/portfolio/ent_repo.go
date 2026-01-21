package portfolio

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entportfolio "github.com/SURF-Innovatie/MORIS/ent/portfolio"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) GetByPersonID(ctx context.Context, personID uuid.UUID) (*entities.Portfolio, error) {
	row, err := r.cli.Portfolio.Query().
		Where(entportfolio.PersonIDEQ(personID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return transform.ToEntityPtr[entities.Portfolio](row), nil
}

func (r *EntRepo) Upsert(ctx context.Context, portfolio entities.Portfolio) (*entities.Portfolio, error) {
	existing, err := r.cli.Portfolio.Query().
		Where(entportfolio.PersonIDEQ(portfolio.PersonID)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if existing == nil {
		row, err := r.cli.Portfolio.Create().
			SetPersonID(portfolio.PersonID).
			SetNillableHeadline(portfolio.Headline).
			SetNillableSummary(portfolio.Summary).
			SetNillableWebsite(portfolio.Website).
			SetShowEmail(portfolio.ShowEmail).
			SetShowOrcid(portfolio.ShowOrcid).
			SetPinnedProjectIds(portfolio.PinnedProjectIDs).
			SetPinnedProductIds(portfolio.PinnedProductIDs).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return transform.ToEntityPtr[entities.Portfolio](row), nil
	}

	row, err := r.cli.Portfolio.UpdateOne(existing).
		SetNillableHeadline(portfolio.Headline).
		SetNillableSummary(portfolio.Summary).
		SetNillableWebsite(portfolio.Website).
		SetShowEmail(portfolio.ShowEmail).
		SetShowOrcid(portfolio.ShowOrcid).
		SetPinnedProjectIds(portfolio.PinnedProjectIDs).
		SetPinnedProductIds(portfolio.PinnedProductIDs).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Portfolio](row), nil
}
