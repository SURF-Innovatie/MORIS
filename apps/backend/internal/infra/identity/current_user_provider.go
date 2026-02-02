package identity

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	authapp "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type CurrentUserProvider struct {
	cli *ent.Client
}

func NewCurrentUserProvider(cli *ent.Client) *CurrentUserProvider {
	return &CurrentUserProvider{cli: cli}
}

func (p *CurrentUserProvider) Current(ctx context.Context) (entities.Principal, error) {
	authUser, ok := httputil.GetUserFromContext(ctx)
	if !ok {
		return entities.Principal{}, fmt.Errorf("no authenticated user in context")
	}

	u, err := p.cli.User.Get(ctx, authUser.User.ID)
	if err != nil {
		return entities.Principal{}, err
	}

	return entities.Principal{
		UserID:     u.ID,
		PersonID:   u.PersonID,
		IsSysAdmin: u.IsSysAdmin,
	}, nil
}

var _ authapp.CurrentUserProvider = (*CurrentUserProvider)(nil)
