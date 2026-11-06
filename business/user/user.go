// Package user provides the core business API for users.
package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/business/sys/database"
	"github.com/ousloob/bookstore/business/user/userdb"
	"golang.org/x/crypto/bcrypt"
)

// Bus manages the set of APIs for user access.
type Bus struct {
	store *userdb.Bus
}

// NewBus constructs a user business API for use.
func NewBus(log *slog.Logger, db *pgxpool.Pool) *Bus {
	return &Bus{
		store: userdb.NewBus(log, db),
	}
}

// Create inserts a new user into the database.
func (b *Bus) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	dbUser := userdb.User{
		ID:           uuid.NewV7().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.store.Create(ctx, dbUser); err != nil {
		if errors.Is(err, database.ErrDuplicatedEntry) {
			return User{}, fmt.Errorf("create: %w", userdb.ErrUniqueEmail)
		}

		return User{}, fmt.Errorf("create: %w", err)
	}

	return User{
		ID:           dbUser.ID,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		Roles:        dbUser.Roles,
		PasswordHash: dbUser.PasswordHash,
		DateCreated:  dbUser.DateCreated,
		DateUpdated:  dbUser.DateUpdated,
	}, nil
}
