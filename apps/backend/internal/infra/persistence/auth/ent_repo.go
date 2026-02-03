package auth

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	userent "github.com/SURF-Innovatie/MORIS/ent/user"
	authapp "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity/readmodels"
	"github.com/google/uuid"
)

type EntRepo struct {
	client  *ent.Client
	userSvc user.Service
}

func NewEntRepo(client *ent.Client, userSvc user.Service) authapp.Repository {
	return &EntRepo{
		client:  client,
		userSvc: userSvc,
	}
}

func (r *EntRepo) GetAccountByEmail(ctx context.Context, email string) (*readmodels.UserAccount, error) {
	return r.userSvc.GetAccountByEmail(ctx, email)
}

func (r *EntRepo) GetAccountByID(ctx context.Context, userID uuid.UUID) (*readmodels.UserAccount, error) {
	return r.userSvc.GetAccount(ctx, userID)
}

func (r *EntRepo) GetPasswordHash(ctx context.Context, userID uuid.UUID) (string, error) {
	u, err := r.client.User.Query().Where(userent.IDEQ(userID)).Only(ctx)
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

var _ authapp.Repository = (*EntRepo)(nil)
