package user

import (
	"fmt"
	"net/mail"
	"time"
	"uuid"

	"github.com/ousloob/bookstore/business/role"
)

// Status is where an account stands. A disabled account is disabled for a
// reason, and the reason decides what happens next: a banned account must
// never come back, an archived one is waiting for its owner to return.
type Status string

// The statuses an account can hold.
const (
	// Active signs in and shops.
	Active Status = "active"

	// Banned was closed by an admin. Only an admin reopens it.
	Banned Status = "banned"

	// Archived went untouched long enough that the data protection rules
	// say to stop keeping it in the open. Signing in brings it back.
	Archived Status = "archived"
)

// ParseStatus turns a string into a Status, rejecting the ones the API does
// not know.
func ParseStatus(s string) (Status, error) {
	switch st := Status(s); st {
	case Active, Banned, Archived:
		return st, nil
	default:
		return "", fmt.Errorf("unknown status %q, want active, banned or archived", s)
	}
}

// User is a user as the business layer sees one.
type User struct {
	ID                uuid.UUID
	Name              string
	Email             mail.Address
	Roles             []role.Role
	PasswordHash      []byte
	Status            Status
	DateCreated       time.Time
	DateUpdated       time.Time
	DateStatusChanged time.Time
	DateLastActive    time.Time
}

// NewUser contains the information needed to create a User.
type NewUser struct {
	Name     string
	Email    mail.Address
	Roles    []role.Role
	Password string
}
