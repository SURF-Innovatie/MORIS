package hierarchy_test

import (
	"context"
	"testing"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/enttest"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/hierarchy"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

func newRBACClient(t *testing.T) *ent.Client {
	t.Helper()
	return enttest.Open(t, "sqlite3", "file:rbacrepo_test?mode=memory&cache=shared&_fk=1")
}

func seedOrgTree(t *testing.T, cli *ent.Client) (rootID, childID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	root, err := cli.OrganisationNode.Create().SetName("root").Save(ctx)
	if err != nil {
		t.Fatalf("create root: %v", err)
	}
	child, err := cli.OrganisationNode.Create().SetName("child").SetParentID(root.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create child: %v", err)
	}

	// closures (including self)
	_, err = cli.OrganisationNodeClosure.Create().SetAncestorID(root.ID).SetDescendantID(root.ID).SetDepth(0).Save(ctx)
	if err != nil {
		t.Fatalf("closure root->root: %v", err)
	}
	_, err = cli.OrganisationNodeClosure.Create().SetAncestorID(child.ID).SetDescendantID(child.ID).SetDepth(0).Save(ctx)
	if err != nil {
		t.Fatalf("closure child->child: %v", err)
	}
	_, err = cli.OrganisationNodeClosure.Create().SetAncestorID(root.ID).SetDescendantID(child.ID).SetDepth(1).Save(ctx)
	if err != nil {
		t.Fatalf("closure root->child: %v", err)
	}

	return root.ID, child.ID
}

func TestAncestorIDs_And_IsAncestor(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := hierarchy.NewEntRepo(cli)
	rootID, childID := seedOrgTree(t, cli)

	ids, err := repo.AncestorIDs(ctx, childID)
	if err != nil {
		t.Fatalf("AncestorIDs: %v", err)
	}

	// Expect both root and child itself (because closures include self)
	if !lo.Contains(ids, rootID) || !lo.Contains(ids, childID) {
		t.Fatalf("expected ancestorIDs to include root and self, got %v", ids)
	}

	ok, err := repo.IsAncestor(ctx, rootID, childID)
	if err != nil {
		t.Fatalf("IsAncestor: %v", err)
	}
	if !ok {
		t.Fatalf("expected root to be ancestor of child")
	}

	ok, err = repo.IsAncestor(ctx, childID, rootID)
	if err != nil {
		t.Fatalf("IsAncestor: %v", err)
	}
	if ok {
		t.Fatalf("did not expect child to be ancestor of root")
	}

	// sanity: closure row exists in DB
	exists, err := cli.OrganisationNodeClosure.Query().
		Where(entclosure.AncestorIDEQ(rootID), entclosure.DescendantIDEQ(childID)).
		Exist(ctx)
	if err != nil {
		t.Fatalf("closure exist: %v", err)
	}
	if !exists {
		t.Fatalf("expected closure row to exist")
	}
}
