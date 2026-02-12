package doi

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/external/doi"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
)

var ErrNotFound = doi.ErrNotFound

type Service interface {
	Resolve(ctx context.Context, doi string) (*dto.Work, error)
}

type service struct {
	client doi.Client
}

func NewService(client doi.Client) Service {
	return &service{client: client}
}

func (s *service) Resolve(ctx context.Context, doiStr string) (*dto.Work, error) {
	w, err := s.client.Resolve(ctx, doiStr)
	if err != nil {
		if errors.Is(err, doi.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := &dto.Work{
		DOI:       w.DOI,
		Title:     w.Title,
		Type:      w.Type,
		Date:      w.Date,
		Publisher: w.Publisher,
	}

	for _, a := range w.Authors {
		out.Authors = append(out.Authors, dto.WorkAuthor{
			Given:  a.Given,
			Family: a.Family,
			Name:   a.Name,
			ORCID:  a.ORCID,
		})
	}

	return out, nil
}
