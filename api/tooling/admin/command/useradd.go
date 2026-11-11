package command

import (
	"bufio"
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
)

// validRoles is the set of roles the API understands. It lives here only
// until M3 gives the business layer a Role type to parse into; then this
// check belongs there, next to the middleware that enforces it. Accepting
// anything in the meantime would create users whose role nothing checks.
var validRoles = []string{"customer", "staff"}

// bcrypt silently ignores everything past the 72nd byte of a password, so
// reject longer ones rather than hand out a credential that is shorter than
// the operator thinks it is.
const maxPasswordBytes = 72

const minPasswordBytes = 8

// UserAdd creates a user. Use it for the first staff account, the one that
// cannot be created through the API because creating staff requires staff.
func UserAdd(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, args []string) error {
	fs := flag.NewFlagSet("useradd", flag.ContinueOnError)
	name := fs.String("name", "", "the person's full name")
	email := fs.String("email", "", "email address, also the login")
	roles := fs.String("roles", "customer", "comma separated roles: customer, staff")

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

func parseRoles(s string) ([]string, error) {
	fields := strings.Split(s, ",")

	parsed := make([]string, 0, len(fields))
	for _, field := range fields {
		role := strings.TrimSpace(field)
		if role == "" {
			continue
		}

		if !slices.Contains(validRoles, role) {
			return nil, fmt.Errorf("unknown role %q, want one of %s", role, strings.Join(validRoles, ", "))
		}

		if slices.Contains(parsed, role) {
			continue
		}

		parsed = append(parsed, role)
	}

	if len(parsed) == 0 {
		return nil, errors.New("at least one role is required")
	}

	return parsed, nil
}

// readPassword takes the password from ADMIN_PASSWORD, so a Kubernetes job
// can supply it from a secret, and otherwise from stdin, so a person can
// pipe it in. There is deliberately no -password flag: an argument is
// visible to anyone who can run ps.
func readPassword() (string, error) {
	pass := os.Getenv("ADMIN_PASSWORD")

	if pass == "" {
		fmt.Fprint(os.Stderr, "Password: ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("reading password: %w", err)
			}

			return "", errors.New("no password given: set ADMIN_PASSWORD or pipe one in")
		}

		pass = strings.TrimSpace(scanner.Text())
		fmt.Fprintln(os.Stderr)
	}

	switch {
	case len(pass) < minPasswordBytes:
		return "", fmt.Errorf("password must be at least %d bytes", minPasswordBytes)
	case len(pass) > maxPasswordBytes:
		return "", fmt.Errorf("password must be at most %d bytes, bcrypt ignores the rest", maxPasswordBytes)
	}

	return pass, nil
}
