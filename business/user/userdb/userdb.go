// Package userdb contains user related CRUD functionality.
package userdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/business/sys/database"
	"github.com/ousloob/bookstore/business/user"
)

// Bus manages the set of APIs for user database access.
type Bus struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

// Check that the storage layer keeps satisfying the business contract.
var _ user.Storer = (*Bus)(nil)

// NewBus constructs a data for api access.
func NewBus(log *slog.Logger, db *pgxpool.Pool) *Bus {
	return &Bus{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (b *Bus) Create(ctx context.Context, usr user.User) error {
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

	dbUsr := toDBUser(usr)

	args := pgx.NamedArgs{
		"user_id":       dbUsr.ID,
		"name":          dbUsr.Name,
		"email":         dbUsr.Email,
		"password_hash": dbUsr.PasswordHash,
		"roles":         dbUsr.Roles,
		"date_created":  dbUsr.DateCreated,
		"date_updated":  dbUsr.DateUpdated,
	}

	if _, err := b.db.Exec(ctx, q, args); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("inserting user: %w", database.ErrDuplicatedEntry)
		}

		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}
