package tx

import "context"

type Manager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
