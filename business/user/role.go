package user

import (
	"fmt"
	"strings"
)

// Role is a set of permissions a user holds. Roles are not ranked: each one
// grants its own area and nothing else, so a user who needs two areas holds
// two roles.
type Role string

// The roles the API understands.
const (
	// RoleCustomer browses, keeps favorites, reserves and buys.
	RoleCustomer Role = "customer"

	// RoleStaff manages the catalog, the stock and the orders.
	RoleStaff Role = "staff"

	// RoleAdmin manages users and their roles. It grants nothing in the
	// store itself: an admin who also runs the shop gives themselves staff.
	RoleAdmin Role = "admin"
)

var roles = []Role{RoleCustomer, RoleStaff, RoleAdmin}

// ParseRole turns a string into a Role, rejecting names the API does not
// know: a user holding an unknown role would hold a role nothing checks.
func ParseRole(s string) (Role, error) {
	for _, r := range roles {
		if string(r) == s {
			return r, nil
		}
	}

	names := make([]string, len(roles))
	for i, r := range roles {
		names[i] = string(r)
	}

	return "", fmt.Errorf("unknown role %q, want one of %s", s, strings.Join(names, ", "))
}
