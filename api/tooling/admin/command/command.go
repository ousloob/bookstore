// Package command implements the subcommands of the admin tool. Every
// command takes the already-open connection pool and the arguments that
// follow its own name, so main.go stays a dispatch table.
package command

import (
	"errors"
	"fmt"
	"io"
)

// ErrHelp is returned when the caller asked for usage rather than for work.
// main.go prints the usage text and exits successfully.
var ErrHelp = errors.New("help requested")

// Usage writes the list of commands.
func Usage(w io.Writer) {
	fmt.Fprint(w, `
Usage: admin <command> [arguments]

Commands:
  migrate              Apply every pending migration
  migrate status       List the migrations and which ones are applied
  migrate down         Roll back the most recently applied migration
  seed                 Load the demo data from migrations/seed
  useradd              Create a user, with the password read from ADMIN_PASSWORD or stdin
  version              Print the build reference
  help                 Print this text

The database connection comes from the environment: DB_USER, DB_PASSWORD,
DB_HOST, DB_NAME, DB_DISABLE_TLS. Run under "infisical run --" to get the
credentials.
`)
}
