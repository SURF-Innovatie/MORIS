package command

import (
	"github.com/SURF-Innovatie/MORIS/ent"
)

type EntClientProvider interface {
	Client() *ent.Client
}
