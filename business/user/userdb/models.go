package userdb

import (
	"time"
	"uuid"

	"github.com/ousloob/bookstore/business/user"
)

// dbUser mirrors the users table with types pgx maps directly.
type dbUser struct {
	ID                uuid.UUID  `db:"user_id"`
	Name              string     `db:"name"`
	Email             string     `db:"email"`
	Roles             []string   `db:"roles"`
	PasswordHash      string     `db:"password_hash"`
	Status            string     `db:"status"`
	DateCreated       time.Time  `db:"date_created"`
	DateUpdated       time.Time  `db:"date_updated"`
	DateStatusChanged time.Time  `db:"date_status_changed"`
	DateLastActive    *time.Time `db:"date_last_active"`
}

func toDBUser(usr user.User) dbUser {
	roles := make([]string, len(usr.Roles))
	for i, r := range usr.Roles {
		roles[i] = string(r)
	}

	// The business layer says "never" with the zero time, the table with NULL.
	var lastActive *time.Time
	if !usr.DateLastActive.IsZero() {
		lastActive = &usr.DateLastActive
	}

	return dbUser{
		ID:                usr.ID,
		Name:              usr.Name,
		Email:             usr.Email.Address,
		Roles:             roles,
		PasswordHash:      string(usr.PasswordHash),
		Status:            string(usr.Status),
		DateCreated:       usr.DateCreated,
		DateUpdated:       usr.DateUpdated,
		DateStatusChanged: usr.DateStatusChanged,
		DateLastActive:    lastActive,
	}
}
