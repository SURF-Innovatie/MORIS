package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
)

type Service interface {
	TotalUserCount(ctx context.Context) (int, error)
}

type service struct {
	client *ent.Client
}

func NewService(client *ent.Client) Service {
	return &service{client: client}
}

func (s *service) TotalUserCount(ctx context.Context) (int, error) {
	return s.client.User.Query().Count(ctx)
}
