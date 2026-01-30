package enttx

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
)

type ctxKey struct{}

type Manager struct {
	cli *ent.Client
}

func NewManager(cli *ent.Client) *Manager {
	return &Manager{cli: cli}
}

func (m *Manager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Nested transaction support: if a tx is already in the context, reuse it.
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.cli.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	ctx = context.WithValue(ctx, ctxKey{}, tx)

	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit()
}

// TxFromContext returns the ent.Tx bound to this context (if any).
func TxFromContext(ctx context.Context) (*ent.Tx, bool) {
	tx, ok := ctx.Value(ctxKey{}).(*ent.Tx)
	return tx, ok
}
