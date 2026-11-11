package userdb

import (
	"time"

	"github.com/ousloob/bookstore/business/user"
)

// dbUser mirrors the users table with types pgx maps directly.
type dbUser struct {
	ID           string    `db:"user_id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Roles        []string  `db:"roles"`
	PasswordHash []byte    `db:"password_hash"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

func toDBUser(usr user.User) dbUser {
	roles := make([]string, len(usr.Roles))
	for i, r := range usr.Roles {
		roles[i] = string(r)
	}

	return dbUser{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		DateCreated:  usr.DateCreated,
		DateUpdated:  usr.DateUpdated,
	}
}
