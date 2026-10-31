# Roadmap

Rough effort estimates are in hours. Order matters: later milestones depend on earlier ones.

## M0. Foundations, cleanup and modernization (about 16 h)

Do this first, while the database code is still tiny. One small commit per item, and update the README's tech stack in the same commit as each change.

**Toolchain and dependencies**

- [x] Rename the Go module from `github.com/ousloob/booksales` to `github.com/ousloob/bookstore`, and fix the README badges
- [x] Move to Go 1.27 (latest patch at the time of writing: 1.27.3) and keep up with patch releases
- [ ] Use `encoding/json/v2` for request and response bodies
- [x] Replace `sqlx` and `lib/pq` with `pgx` v5 (`pgxpool`, named arguments; use `CollectRows` and `RowToStructByName` once there are read queries)
- [x] Update all remaining dependencies, and keep `govulncheck` green
- [x] Drop the `vendor/` directory and rely on the module cache (`go.sum` already pins dependencies)
- [x] Replace `github.com/google/uuid` with the standard library `uuid` package (Go 1.27+)
- [x] Evaluate replacing `golang.org/x/crypto/bcrypt`: the standard library has no bcrypt or argon2id (only `crypto/pbkdf2`), and `x/crypto` is maintained by the Go team with monthly releases, so keep it
- [ ] Evaluate replacing `chi` with the standard library router (pattern routing, Go 1.22+)

**CI and tooling**

- [x] Replace CircleCI with GitHub Actions: tests, `go vet`, `staticcheck`, `govulncheck`
- [x] Pin tool versions instead of using `@latest`, and read the Go version from `go.mod`
- [ ] Enable Dependabot for Go modules and for the workflow files

**Database**

- [ ] Run PostgreSQL 18 (current stable release) in development and in CI; consider 19 after its final release
- [ ] Add schema migrations (Goose) for users and books
- [ ] Map PostgreSQL unique-violation errors to `ErrDBDuplicatedEntry`

**Bugs and safety**

- [x] Fix the user insert query (`date_updated` was missing its `:` in the `VALUES` list; fixed when the query moved to `pgx` named arguments)
- [ ] Fix error handling in `userbiz.Create`: errors other than a duplicate email are currently swallowed
- [ ] Stop logging the full configuration at startup (it includes the database password), and remove the default database user and password from the code
- [ ] Add `.env.example` and document the environment variables
- [ ] Pass the logger through the routing setup (it is currently nil)

**Tests**

- [ ] Replace the placeholder tests with real ones, using a PostgreSQL container for integration tests

## M1. Catalogue (about 8 h)

- [ ] Book model: ISBN as a string with validation (ISBN-10 can end in `X`), price in cents, stock, formats (physical, digital)
- [ ] Staff endpoints: create, update and delete books
- [ ] Public endpoints: list, search by title, author or ISBN, pagination and sorting
- [ ] Consistent JSON error responses and request validation
- [ ] OpenAPI description of the API

## M2. Accounts and authentication (about 8 h)

- [ ] Register and log in with email and password (bcrypt, or argon2id)
- [ ] Short-lived access tokens, validated on every request (signature, issuer, audience, expiry)
- [ ] Roles: customer and staff, enforced per route
- [ ] Rate limiting on login
- [ ] Later, optional: sign in with an external identity provider (OpenID Connect)

## M3. Favorites, reservations and orders (about 10 h)

- [ ] Favorites: add, remove and list
- [ ] Reservations: hold a copy for a set number of days, release it automatically when it expires
- [ ] Orders for physical copies, with stock decremented safely inside a database transaction (no overselling under concurrent requests)
- [ ] Order status: created, ready, collected, cancelled

## M4. Digital books (about 6 h)

- [ ] Digital library per customer
- [ ] File storage behind an interface (local disk first)
- [ ] Download through time-limited signed links

## M5. Store management ideas (pick what is useful)

- [ ] Inventory adjustments with an audit log (who changed what, and when)
- [ ] Low-stock alerts
- [ ] Back-in-stock notification list
- [ ] Sales and bestseller reports
- [ ] Book metadata import by ISBN from a public source such as Open Library
- [ ] CSV import and export of the catalogue
- [ ] Returns and cancellations
- [ ] Discount codes
- [ ] Reviews and ratings
- [ ] Events (author signings, reading groups)

## M6. Operations and quality (about 4 h)

- [ ] Dockerfile and `docker-compose.yml` for local development
- [ ] Health and readiness endpoints
- [ ] Metrics and request tracing
- [ ] Optional: Kubernetes manifests
