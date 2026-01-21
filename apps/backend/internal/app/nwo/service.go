package nwo

import (
	"context"
	"errors"

	ex "github.com/SURF-Innovatie/MORIS/external/nwo"
)

var (
	ErrNotFound = errors.New("nwo_not_found")
)

// Service provides business logic for NWO Open API operations
type Service interface {
	// GetProjects queries NWO projects with optional filters
	GetProjects(ctx context.Context, opts *ex.QueryOptions) (*ex.ProjectsResponse, error)
	// GetProject retrieves a single project by project_id
	GetProject(ctx context.Context, projectID string) (*ex.Project, error)
}

type service struct {
	client ex.Client
}

// NewService creates a new NWO service
func NewService(client ex.Client) Service {
	return &service{client: client}
}

func (s *service) GetProjects(ctx context.Context, opts *ex.QueryOptions) (*ex.ProjectsResponse, error) {
	resp, err := s.client.GetProjects(ctx, opts)
	if errors.Is(err, ex.ErrNotFound) {
		return nil, ErrNotFound
	}
	return resp, err
}

func (s *service) GetProject(ctx context.Context, projectID string) (*ex.Project, error) {
	project, err := s.client.GetProject(ctx, projectID)
	if errors.Is(err, ex.ErrNotFound) {
		return nil, ErrNotFound
	}
	return project, err
}
