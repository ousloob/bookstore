// Package user provides the core business API for users.
package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"uuid"

	"github.com/ousloob/bookstore/foundation/database"
	"golang.org/x/crypto/bcrypt"
)

// ErrUniqueEmail is returned when the email is already taken.
var ErrUniqueEmail = errors.New("email is not unique")

// Storer is the storage this business API needs. userdb implements it.
type Storer interface {
	Create(ctx context.Context, usr User) error
}

// Bus manages the set of APIs for user access.
type Bus struct {
	log   *slog.Logger
	store Storer
}

// NewBus constructs a user business API for use.
func NewBus(log *slog.Logger, store Storer) *Bus {
	return &Bus{
		log:   log,
		store: store,
	}
}

// Create inserts a new user into the database.
func (b *Bus) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	// Providers treat the local part as case-insensitive, so one spelling
	// keeps Bob@x.com and bob@x.com from becoming two accounts.
	email := nu.Email
	email.Address = strings.ToLower(email.Address)

	usr := User{
		ID:           uuid.NewV7(),
		Name:         nu.Name,
		Email:        email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		Status:       Active,

		DateCreated:       now,
		DateUpdated:       now,
		DateStatusChanged: now,
	}

	if err := b.store.Create(ctx, usr); err != nil {
		if errors.Is(err, database.ErrDuplicatedEntry) {
			return User{}, fmt.Errorf("create: %w", ErrUniqueEmail)
		}

		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
