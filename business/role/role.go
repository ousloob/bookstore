// Package role holds the roles a user can hold. It imports nothing, so every
// other layer is free to depend on it: the business packages that assign
// roles, the middleware that checks them and the token that carries them.
package role

import (
	"fmt"
	"strings"
)

// Role is one area of the system a user is allowed to work in. Roles are not
// ranked and do not contain one another: each grants its own area and nothing
// else, so a user who needs two areas holds two roles.
type Role string

// The roles the API understands.
const (
	// Customer browses, keeps favorites, reserves and buys.
	Customer Role = "customer"

	// Staff manages the catalog, the stock and the orders.
	Staff Role = "staff"

	// Admin manages users and their roles. It grants nothing in the store
	// itself: an admin who also runs the shop holds Staff as well.
	Admin Role = "admin"
)

// All is every role, in no meaningful order.
var All = []Role{Customer, Staff, Admin}

// Parse turns a string into a Role, rejecting names the API does not know: a
// user holding an unknown role would hold a role that nothing ever checks.
func Parse(s string) (Role, error) {
	for _, r := range All {
		if string(r) == s {
			return r, nil
		}
	}

	names := make([]string, len(All))
	for i, r := range All {
		names[i] = string(r)
	}

	return "", fmt.Errorf("unknown role %q, want one of %s", s, strings.Join(names, ", "))
}
