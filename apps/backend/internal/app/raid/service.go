package raid

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/google/uuid"
)

type Service interface {
	MintRaid(ctx context.Context, userID uuid.UUID, req *raid.RAiDCreateRequest) (*raid.RAiDDto, error)
	UpdateRaid(ctx context.Context, userID uuid.UUID, prefix, suffix string, req *raid.RAiDUpdateRequest) (*raid.RAiDDto, error)
	FindRaid(ctx context.Context, userID uuid.UUID, prefix, suffix string) (*raid.RAiDDto, error)
	FindAllRaids(ctx context.Context, userID uuid.UUID) ([]raid.RAiDDto, error)
}

type Client interface {
	MintRaid(ctx context.Context, req *raid.RAiDCreateRequest) (*raid.RAiDDto, error)
	UpdateRaid(ctx context.Context, prefix, suffix string, req *raid.RAiDUpdateRequest) (*raid.RAiDDto, error)
	FindRaid(ctx context.Context, prefix, suffix string) (*raid.RAiDDto, error)
	FindAllRaids(ctx context.Context) ([]raid.RAiDDto, error)
}

type service struct {
	client Client
}

func NewService(client Client) Service {
	return &service{client: client}
}

func (s *service) MintRaid(ctx context.Context, userID uuid.UUID, req *raid.RAiDCreateRequest) (*raid.RAiDDto, error) {
	// TODO: Use userID for auditing or validation if needed
	return s.client.MintRaid(ctx, req)
}

func (s *service) UpdateRaid(ctx context.Context, userID uuid.UUID, prefix, suffix string, req *raid.RAiDUpdateRequest) (*raid.RAiDDto, error) {
	return s.client.UpdateRaid(ctx, prefix, suffix, req)
}

func (s *service) FindRaid(ctx context.Context, userID uuid.UUID, prefix, suffix string) (*raid.RAiDDto, error) {
	return s.client.FindRaid(ctx, prefix, suffix)
}

func (s *service) FindAllRaids(ctx context.Context, userID uuid.UUID) ([]raid.RAiDDto, error) {
	return s.client.FindAllRaids(ctx)
}
