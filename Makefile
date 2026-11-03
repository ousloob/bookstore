SHELL := /bin/sh

# Keep these in sync with .github/workflows/ci.yml.
STATICCHECK := honnef.co/go/tools/cmd/staticcheck@v0.8.2
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.9.0

# ==============================================================================
# Module support

.PHONY: tidy
tidy:
	go mod tidy

# ==============================================================================
# Local PostgreSQL
#
# DB_USER and DB_PASSWORD come from Infisical, so these targets are run as
# `infisical run -- make db-up`. The credentials are baked into the data
# directory the first time the server initializes it: changing them in
# Infisical afterwards has no effect until `make db-reset`.

POSTGRES_IMAGE := postgres:18-alpine
DB_CONTAINER   := bookstore-db
DB_VOLUME      := bookstore-db-data

.PHONY: db-up db-down db-reset db-psql
db-up:
	docker run -d --name $(DB_CONTAINER) \
		-e POSTGRES_USER="$${DB_USER:?run via: infisical run -- make db-up}" \
		-e POSTGRES_PASSWORD="$${DB_PASSWORD:?run via: infisical run -- make db-up}" \
		-e POSTGRES_DB=bookdb \
		-p 5432:5432 \
		-v $(DB_VOLUME):/var/lib/postgresql \
		$(POSTGRES_IMAGE)

db-down:
	docker rm -f $(DB_CONTAINER)

# Destructive: drops the database and every row in it.
db-reset: db-down
	docker volume rm $(DB_VOLUME)

db-psql:
	docker exec -it $(DB_CONTAINER) psql -U "$${DB_USER:?}" -d bookdb

# ==============================================================================
# Running the API locally

.PHONY: run
run:
	go run api/sales/main.go

# ==============================================================================
# Code quality check locally

.PHONY: test lint vulncheck fix check

# Run all checks
check: lint vulncheck test

lint:
	CGO_ENABLED=0 go vet ./...
	go fix -diff ./...
	go run $(STATICCHECK) -checks=all ./...

vulncheck:
	go run $(GOVULNCHECK) ./...

test:
	go test -count=1 ./...

# Apply the modernizations that `make lint` reports
fix:
	go fix ./...

# ==============================================================================
