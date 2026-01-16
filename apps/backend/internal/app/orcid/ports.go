package orcid

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/external/orcid"
)

type OrcidClient interface {
	AuthURL() (string, error)
	ExchangeCode(ctx context.Context, code string) (string, error)
	SearchExpanded(ctx context.Context, query string) ([]orcid.OrcidPerson, error)
}
