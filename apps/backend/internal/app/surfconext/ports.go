package surfconext

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/external/surfconext"
)

// Client defines the interface for interacting with a SURFconext OIDC provider.
type Client interface {
	AuthURL(ctx context.Context) (string, error)
	ExchangeCode(ctx context.Context, code string) (*surfconext.Claims, error)
}
