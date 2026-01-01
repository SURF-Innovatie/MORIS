package bootstrap

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/projectrole"
)

type RoleDefaults struct {
	Repo projectrole.Repository
}

func (b RoleDefaults) Ensure(ctx context.Context) error {
	defs := []struct{ Key, Name string }{
		{Key: "contributor", Name: "Contributor"},
		{Key: "lead", Name: "Project lead"},
	}
	for _, d := range defs {
		if err := b.Repo.Upsert(ctx, d.Key, d.Name); err != nil {
			return err
		}
	}
	return nil
}
