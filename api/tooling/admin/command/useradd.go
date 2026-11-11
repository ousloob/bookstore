package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/business/user"
	"github.com/ousloob/bookstore/business/user/userdb"
	"golang.org/x/term"
)

// bcrypt silently ignores everything past the 72nd byte of a password, so
// reject longer ones rather than hand out a credential that is shorter than
// the operator thinks it is.
const maxPasswordBytes = 72

const minPasswordBytes = 8

// UserAdd creates a user. It exists for the first admin, the one account the
// API cannot create because creating users with roles takes an admin. Run it
// by hand, once: the password is typed at a hidden prompt and is never stored
// anywhere but as a hash.
func UserAdd(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, args []string) error {
	fs := flag.NewFlagSet("useradd", flag.ContinueOnError)
	name := fs.String("name", "", "the person's full name")
	email := fs.String("email", "", "email address, also the login")
	roles := fs.String("roles", "admin", "comma separated roles: customer, staff, admin")

	if err := fs.Parse(args); err != nil {
		return ErrHelp
	}

	if *name == "" || *email == "" {
		return errors.New("useradd: -name and -email are required")
	}

	addr, err := mail.ParseAddress(*email)
	if err != nil {
		return fmt.Errorf("useradd: parsing email: %w", err)
	}

	parsedRoles, err := parseRoles(*roles)
	if err != nil {
		return fmt.Errorf("useradd: %w", err)
	}

	pass, err := readPassword()
	if err != nil {
		return fmt.Errorf("useradd: %w", err)
	}

	usrBus := user.NewBus(log, userdb.NewBus(log, db))

	usr, err := usrBus.Create(ctx, user.NewUser{
		Name:     *name,
		Email:    *addr,
		Roles:    parsedRoles,
		Password: pass,
	}, time.Now().UTC())
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return fmt.Errorf("useradd: %s already exists", addr.Address)
		}

		return fmt.Errorf("useradd: %w", err)
	}

	log.InfoContext(ctx, "useradd", "status", "user created", "id", usr.ID, "email", usr.Email.Address, "roles", usr.Roles)

	return nil
}

func parseRoles(s string) ([]user.Role, error) {
	var parsed []user.Role

	for field := range strings.SplitSeq(s, ",") {
		name := strings.TrimSpace(field)
		if name == "" {
			continue
		}

		role, err := user.ParseRole(name)
		if err != nil {
			return nil, err
		}

		if !slices.Contains(parsed, role) {
			parsed = append(parsed, role)
		}
	}

	if len(parsed) == 0 {
		return nil, errors.New("at least one role is required")
	}

	return parsed, nil
}

// readPassword prompts twice on the terminal with echo off. It refuses to
// read from a pipe or a file: a password that arrives that way has been
// written down somewhere first, and this command is meant to be run by a
// person.
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", errors.New("must be run from a terminal, it prompts for the password")
	}

	pass, err := prompt(fd, "Password: ")
	if err != nil {
		return "", err
	}

	switch {
	case len(pass) < minPasswordBytes:
		return "", fmt.Errorf("password must be at least %d bytes", minPasswordBytes)
	case len(pass) > maxPasswordBytes:
		return "", fmt.Errorf("password must be at most %d bytes, bcrypt ignores the rest", maxPasswordBytes)
	}

	again, err := prompt(fd, "Again: ")
	if err != nil {
		return "", err
	}

	if pass != again {
		return "", errors.New("the passwords do not match")
	}

	return pass, nil
}

func prompt(fd int, label string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	b, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("reading password: %w", err)
	}

	return string(b), nil
}
