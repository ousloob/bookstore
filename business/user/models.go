package user

import (
	"net/mail"
	"time"
	"uuid"

	"github.com/ousloob/bookstore/business/role"
)

// User is a user as the business layer sees one.
type User struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	Roles        []role.Role
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contains the information needed to create a User.
type NewUser struct {
	Name     string
	Email    mail.Address
	Roles    []role.Role
	Password string
}
