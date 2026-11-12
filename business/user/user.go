// Package user provides the core business API for users.
package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

	usr := User{
		ID:           uuid.NewV7(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.store.Create(ctx, usr); err != nil {
		if errors.Is(err, database.ErrDuplicatedEntry) {
			return User{}, fmt.Errorf("create: %w", ErrUniqueEmail)
		}

		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
