package identity

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	authapp "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type CurrentUserProvider struct {
	cli *ent.Client
}

func NewCurrentUserProvider(cli *ent.Client) *CurrentUserProvider {
	return &CurrentUserProvider{cli: cli}
}

func (p *CurrentUserProvider) Current(ctx context.Context) (identity.Principal, error) {
	authUser, ok := httputil.GetUserFromContext(ctx)
	if !ok {
		return identity.Principal{}, fmt.Errorf("no authenticated user in context")
	}

	u, err := p.cli.User.Get(ctx, authUser.User.ID)
	if err != nil {
		return identity.Principal{}, err
	}

	return identity.Principal{
		UserID:     u.ID,
		PersonID:   u.PersonID,
		IsSysAdmin: u.IsSysAdmin,
	}, nil
}

var _ authapp.CurrentUserProvider = (*CurrentUserProvider)(nil)
