// Package userdb contains user related CRUD functionality.
package userdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUniqueEmail = errors.New("email is not unique")

// Bus manages the set of APIs for user access.
type Bus struct {
	log *log.Logger
	db  *pgxpool.Pool
}

// NewBus constructs a data for api access.
func NewBus(log *log.Logger, db *pgxpool.Pool) *Bus {
	return &Bus{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (b *Bus) Create(ctx context.Context, usr User) error {
	const q = `
	INSERT INTO users (
        user_id,
        name,
        email,
        password_hash,
        roles,
        date_created,
        date_updated
    )
	VALUES (
        @user_id,
        @name,
        @email,
        @password_hash,
        @roles,
        @date_created,
        @date_updated
    )`

	args := pgx.NamedArgs{
		"user_id":       usr.ID,
		"name":          usr.Name,
		"email":         usr.Email,
		"password_hash": usr.PasswordHash,
		"roles":         usr.Roles,
		"date_created":  usr.DateCreated,
		"date_updated":  usr.DateUpdated,
	}

	if _, err := b.db.Exec(ctx, q, args); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}
