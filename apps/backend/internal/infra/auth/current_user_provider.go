package auth

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
)

type currentUser struct {
	userID   uuid.UUID
	personID uuid.UUID
}

func (u currentUser) UserID() uuid.UUID   { return u.userID }
func (u currentUser) PersonID() uuid.UUID { return u.personID }

type CurrentUserProvider struct {
	cli *ent.Client
}

func NewCurrentUserProvider(cli *ent.Client) *CurrentUserProvider {
	return &CurrentUserProvider{cli: cli}
}

func (p *CurrentUserProvider) Current(ctx context.Context) (appauth.CurrentUser, error) {
	authUser, ok := httputil.GetUserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no authenticated user in context")
	}

	u, err := p.cli.User.Get(ctx, authUser.User.ID)
	if err != nil {
		return nil, err
	}

	return currentUser{userID: u.ID, personID: u.PersonID}, nil
}
