package organisation_test

import (
	"context"
	"testing"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/enttest"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	_ "github.com/mattn/go-sqlite3"
)

func newOrgTestClient(t *testing.T) *ent.Client {
	t.Helper()
	return enttest.Open(t, "sqlite3", "file:orgrepo_test?mode=memory&cache=shared&_fk=1")
}

func TestOrganisationRepo_UsesTxFromContext_Rollback(t *testing.T) {
	cli := newOrgTestClient(t)
	defer cli.Close()

	repo := organisation.NewEntRepo(cli)
	txMgr := enttx.NewManager(cli)

	err := txMgr.WithTx(context.Background(), func(ctx context.Context) error {
		_, err := repo.CreateNode(ctx, "root", nil, nil, nil, nil, "root")
		if err != nil {
			t.Fatalf("CreateNode failed: %v", err)
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

func TestOrganisationRepo_UsesTxFromContext_Commit(t *testing.T) {
	cli := newOrgTestClient(t)
	defer cli.Close()

	repo := organisation.NewEntRepo(cli)
	txMgr := enttx.NewManager(cli)

	err := txMgr.WithTx(context.Background(), func(ctx context.Context) error {
		_, err := repo.CreateNode(ctx, "root", nil, nil, nil, nil, "root")
		return err
	})
	if err != nil {
		t.Fatalf("WithTx failed: %v", err)
	}

	n, err := cli.OrganisationNode.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 node after commit, got %d", n)
	}
}
