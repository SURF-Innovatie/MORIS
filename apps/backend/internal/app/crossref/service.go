package crossref

import (
	"context"
	"errors"

	ex "github.com/SURF-Innovatie/MORIS/external/crossref"
)

var (
	ErrNotFound = errors.New("crossref_not_found")
)

type Service interface {
	GetWork(ctx context.Context, doi string) (*ex.Work, error)
	GetWorks(ctx context.Context, query string, limit int) ([]ex.Work, error)
	GetJournal(ctx context.Context, issn string) (*ex.Journal, error)
	GetJournals(ctx context.Context, query string, limit int) ([]ex.Journal, error)
}

type service struct {
	client ex.Client
}

func NewService(client ex.Client) Service {
	return &service{client: client}
}

func (s *service) GetWork(ctx context.Context, doi string) (*ex.Work, error) {
	w, err := s.client.GetWork(ctx, doi)
	if errors.Is(err, ex.ErrNotFound) {
		return nil, ErrNotFound
	}
	return w, err
}

func (s *service) GetWorks(ctx context.Context, query string, limit int) ([]ex.Work, error) {
	ws, err := s.client.GetWorks(ctx, query, limit)
	if errors.Is(err, ex.ErrNotFound) {
		return nil, ErrNotFound
	}
	return ws, err
}

func (s *service) GetJournal(ctx context.Context, issn string) (*ex.Journal, error) {
	j, err := s.client.GetJournal(ctx, issn)
	if errors.Is(err, ex.ErrNotFound) {
		return nil, ErrNotFound
	}
	return j, err
}

func (s *service) GetJournals(ctx context.Context, query string, limit int) ([]ex.Journal, error) {
	js, err := s.client.GetJournals(ctx, query, limit)
	if errors.Is(err, ex.ErrNotFound) {
		return nil, ErrNotFound
	}
	return js, err
}
