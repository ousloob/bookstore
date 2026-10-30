SHELL := /bin/sh

# Keep these in sync with .github/workflows/ci.yml.
STATICCHECK := honnef.co/go/tools/cmd/staticcheck@v0.8.2
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.9.0

# ==============================================================================
# Module support

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

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
