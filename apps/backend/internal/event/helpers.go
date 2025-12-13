package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/google/uuid"
)

func ResolveUser(ctx context.Context, cli *ent.Client, creatorID uuid.UUID) (*ent.User, error) {
	if creatorID == uuid.Nil {
		return nil, nil
	}

	// Try as User.ID first
	u, err := cli.User.Get(ctx, creatorID)
	if err == nil {
		return u, nil
	}

	// Fallback: treat creatorID as Person.ID
	u, err = cli.User.Query().Where(user.PersonIDEQ(creatorID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}
