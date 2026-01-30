package enttx_test

import (
	"context"
	"testing"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/enttest"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	_ "github.com/mattn/go-sqlite3"
)

func newTestClient(t *testing.T) *ent.Client {
	t.Helper()
	return enttest.Open(t, "sqlite3", "file:enttx_test?mode=memory&cache=shared&_fk=1")
}

func TestManager_WithTx_Commits(t *testing.T) {
	cli := newTestClient(t)
	defer cli.Close()

	m := enttx.NewManager(cli)

	err := m.WithTx(context.Background(), func(ctx context.Context) error {
		tx, ok := enttx.TxFromContext(ctx)
		if !ok || tx == nil {
			t.Fatalf("expected tx in context")
		}
		_, err := tx.OrganisationNode.Create().SetName("root").Save(ctx)
		return err
	})
	if err != nil {
		t.Fatalf("WithTx returned error: %v", err)
	}

	n, err := cli.OrganisationNode.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 node, got %d", n)
	}
}

func TestManager_WithTx_RollsBackOnError(t *testing.T) {
	cli := newTestClient(t)
	defer cli.Close()

	m := enttx.NewManager(cli)

	err := m.WithTx(context.Background(), func(ctx context.Context) error {
		tx, _ := enttx.TxFromContext(ctx)
		_, err := tx.OrganisationNode.Create().SetName("root").Save(ctx)
		if err != nil {
			return err
		}
		return context.Canceled
	})
	if err == nil {
		t.Fatalf("expected error")
	}

	n, err := cli.OrganisationNode.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 nodes after rollback, got %d", n)
	}
}

func TestManager_WithTx_NestedReusesSameTx(t *testing.T) {
	cli := newTestClient(t)
	defer cli.Close()

	m := enttx.NewManager(cli)

	err := m.WithTx(context.Background(), func(ctx context.Context) error {
		tx1, ok := enttx.TxFromContext(ctx)
		if !ok || tx1 == nil {
			t.Fatalf("expected tx1 in context")
		}

		return m.WithTx(ctx, func(ctx2 context.Context) error {
			tx2, ok := enttx.TxFromContext(ctx2)
			if !ok || tx2 == nil {
				t.Fatalf("expected tx2 in context")
			}
			if tx1 != tx2 {
				t.Fatalf("expected nested WithTx to reuse same tx pointer")
			}
			return nil
		})
	})
	if err != nil {
		t.Fatalf("WithTx returned error: %v", err)
	}
}
