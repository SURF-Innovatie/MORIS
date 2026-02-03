package readmodels

import "github.com/SURF-Innovatie/MORIS/internal/domain/identity"

type UserAccount struct {
	User   identity.User
	Person identity.Person
}
