# Roadmap

Working document: finished items are checked off. Milestones are in dependency order. One small commit per item; update the README's tech stack in the same commit.

## M0. Groundwork

**Local environment**

- [x] [Podman](https://podman.io) as the container engine, [kind](https://kind.sigs.k8s.io) on its Podman provider; kind only, no Docker Compose
- [x] `make dev-up` creates the cluster, builds and loads the image and applies `zarf/k8s/dev`; `make dev-down` deletes it
- [x] PostgreSQL 18 in the Kustomize dev overlay `zarf/k8s/dev` (StatefulSet, volume, Service), credentials from the same `DB_USER` and `DB_PASSWORD` the app reads; `make dev-db-forward` port-forwards it so `make run` works against a real database

**Configuration and secrets**

- [ ] Secrets (database user and password, token signing key, anything sensitive) live in [Infisical](https://infisical.com) and reach the app as environment variables: `infisical run` locally, the Infisical operator in Kubernetes

**Database**

- [x] Expose every `pgxpool.Config` setting (`MaxConns`, `MinConns`, `MinIdleConns`, `MaxConnLifetime`, `MaxConnLifetimeJitter`, `MaxConnIdleTime`, `HealthCheckPeriod`, `PingTimeout`); unset keeps the pgx default
- [x] [Goose](https://github.com/pressly/goose) migrations: formatted SQL files in `migrations/` at the repository root, embedded with `go:embed`
- [x] `api/tooling/admin` CLI: `migrate`, `migrate status`, `migrate down`, `seed`, `useradd`; same `DB_*` configuration and same image as the API
- [x] Stop migrating at API startup: the admin tool owns the schema
- [ ] Separate database roles: a migrator role with DDL rights for the admin tool, a DML-only role for the API
- [x] Seed data for demos: idempotent SQL in `migrations/seed/`, run by `make seed`
- [x] Map unique violations (`23505`) to `ErrDuplicatedEntry`
- [x] UUID v7 IDs with the standard library `uuid.NewV7()`; v4 or raw `crypto/rand` bytes for anything that must be unguessable
- [ ] `database.WithTx` helper for multi-step writes

**Code cleanup**

- [x] Fix `userbiz.Create`: every error except a duplicate email is swallowed
- [x] Business and storage layers take `*slog.Logger` instead of `*log.Logger`
- [ ] `encoding/json/v2` for request and response bodies
- [x] Optional: rename `userbiz` and `bookbiz` to `business/user` and `business/book` (storage stays `userdb`, `bookdb`)

**Tooling and CI**

- [ ] [golangci-lint v2](https://golangci-lint.run) with `gosec` enabled (minus G104 and G115, which duplicate errcheck and flag every integer conversion) replaces `staticcheck`; keep `go fix -diff` and `govulncheck`
- [ ] Dependabot for Go modules and workflow files
- [ ] [GoReleaser](https://goreleaser.com) on tag push: binaries, checksums, changelog (container image added in M8)
- [ ] Integration tests against PostgreSQL 18 with [testcontainers-go](https://golang.testcontainers.org), run in CI; replace the placeholder tests

## M1. API contract

- [ ] All routes under `/v1`
- [ ] Every error as [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) problem details (`application/problem+json`), including the mux's 404 and 405, with an `errors` extension for field validation
- [ ] List conventions: `page`, `page_size` (default 20, max 100), `sort=field` or `sort=-field` from a per-resource allow-list, filters as query parameters
- [ ] Hand-written OpenAPI 3.1 in `api/openapi.yaml`, no code generation, served at `GET /v1/openapi.yaml`, linted in CI

## M2. Users and books

- [ ] Users: create, get, update, delete; list with filters (`email`, `role`), sorting and pagination; the password hash never leaves the storage layer
- [ ] Book model: ISBN string with validation, price in cents, stock, formats (physical, digital)
- [ ] Books: create, update, delete (staff); get and list with filters (`title`, `author`, `isbn`, `format`, `in_stock`), sorting and pagination

## M3. Authentication and roles

- [ ] Passwords hashed with argon2id (`golang.org/x/crypto/argon2`), PHC format, rehash on login when parameters change
- [ ] Login with email and password, same error whether or not the email exists
- [ ] Access token: short-lived JWT, Ed25519, `golang-jwt/jwt/v5`, validated on every request
- [ ] Refresh token: random bytes from `crypto/rand`, stored hashed, bound to a session, rotated on every use; reuse ends the session; logout revokes the session
- [ ] `Authorizer` interface behind one middleware; first implementation: roles `customer` and `staff` from the token. Move to [Casbin](https://casbin.org) when roles multiply
- [ ] Rate limiting on login and refresh
- [ ] Security audit log: login, failed login, refresh reuse, logout

## M4. Favorites, reservations and orders

- [ ] Favorites: add, remove, list
- [ ] Reservations: hold a copy for a number of days, auto-release on expiry
- [ ] Orders with status: created, ready, collected, cancelled
- [ ] Order creation in one transaction: conditional stock decrement (`stock >= qty`, zero rows means out of stock), then order and lines; `CHECK (stock >= 0)`
- [ ] Idempotency keys: `Idempotency-Key` header stored with the user, a body hash and the response, in the same transaction; replay returns the stored response, different body 422, in flight 409, expiry 24 h
- [ ] Tests: concurrent orders for the last copy; same request twice with one key

## M5. Payments

- [ ] Payment providers behind one interface, so every method shares the same order flow
- [ ] Card payments through the provider's hosted checkout; no card data reaches this API
- [ ] [Wero](https://wero-wallet.eu) account-to-account payments through a PSP that supports it (merchant acceptance is still rolling out, so availability depends on country)
- [ ] Order states for payment: pending, paid, failed, refunded; stock stays held while payment is pending and is released on timeout
- [ ] Provider webhook endpoint: signature verified, idempotent, replay-safe
- [ ] Refunds and cancellations
- [ ] Tests against the provider sandbox

## M6. Digital books

- [ ] Digital library per customer
- [ ] File storage behind an interface, local disk first
- [ ] Download through time-limited signed links

## M7. Observability

- [ ] `log/slog`: `tint` in development, JSON in production; request ID and trace ID on every line
- [ ] OpenTelemetry: `otelhttp`, `otelpgx`, request metrics, OTLP export, set up in `support/telemetry`
- [x] `/livez` and `/readyz` (pings the database)
- [ ] Grafana, Prometheus, Tempo, Loki and Alloy deployed in the kind dev overlay (M8), one starter dashboard committed

## M8. Release and deployment

- [ ] Containerfile (built in M0 at `zarf/containers`: multi-stage, alpine, non-root): image pushed to GHCR by GoReleaser
- [ ] Kustomize: `zarf/k8s/base` (done in M0: Deployment, Service; still to come: ConfigMap, migration Job) and a `zarf/k8s/prod` overlay; the dev overlay from M0 adds the app, the observability stack and an `InfisicalSecret`
- [ ] CI job: build the image, load into kind, apply dev, check `/readyz`

## M9. Maybe later

- [ ] Argo CD and a Helm chart, once there is more than one environment
- [ ] Web front end, once payments work end to end
- [ ] Store tools: inventory audit log, low-stock alerts, sales reports, ISBN import from Open Library, CSV import and export
- [ ] Returns, discount codes, reviews, events
