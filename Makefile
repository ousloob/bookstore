SHELL := /bin/sh

# Keep these in sync with .github/workflows/ci.yml.
STATICCHECK := honnef.co/go/tools/cmd/staticcheck@v0.8.2
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.9.0

# ==============================================================================
# Help

.DEFAULT_GOAL := help

CYAN  := \033[36m
RESET := \033[0m

# Every target documents itself with a `## description` comment on its own
# line. This prints one row per such target: awk splits each line on the
# `:` and the `## `, leaving the target name in $$1 and the description in $$2.
.PHONY: help
help: ## Show this help
	@awk 'BEGIN { FS = ":.*## " } \
		/^[a-zA-Z_-]+:.*## / { printf "  $(CYAN)%-14s$(RESET) %s\n", $$1, $$2 }' \
		$(MAKEFILE_LIST)

# ==============================================================================
# Dev environment (kind + podman)

KIND_CLUSTER := bookstore-cluster
SALES_IMAGE  := localhost/bookstore/sales-api:dev
NAMESPACE    := bookstore-system
BUILD_REF    := $(shell git rev-parse --short HEAD)

# kind's podman support is experimental and must be opted into.
export KIND_EXPERIMENTAL_PROVIDER := podman

.PHONY: dev-up dev-down dev-build dev-load dev-secrets dev-apply dev-update dev-status dev-logs

dev-up: ## Create the kind cluster and deploy everything into it
	kind create cluster --name $(KIND_CLUSTER) --config zarf/k8s/dev/kind-config.yaml
	$(MAKE) dev-update

dev-down: ## Delete the kind cluster
	kind delete cluster --name $(KIND_CLUSTER)

dev-build: ## Build the sales-api image with podman
	podman build -f zarf/containers/Containerfile.sales -t $(SALES_IMAGE) --build-arg BUILD_REF=$(BUILD_REF) .

dev-load: ## Load the sales-api image into the kind cluster
	podman save $(SALES_IMAGE) -o /tmp/sales-api.tar
	kind load image-archive /tmp/sales-api.tar --name $(KIND_CLUSTER)
	rm /tmp/sales-api.tar

dev-secrets: ## Sync the database credentials from Infisical into the cluster
	@test -f .infisical.json || { \
		echo "No .infisical.json found. Set up Infisical first:"; \
		echo "    infisical login"; \
		echo "    infisical init"; \
		exit 1; \
	}
	kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	infisical run -- sh -c 'kubectl create secret generic database-credentials \
		--namespace=$(NAMESPACE) \
		--from-literal=DB_USER="$$DB_USER" \
		--from-literal=DB_PASSWORD="$$DB_PASSWORD" \
		--dry-run=client -o yaml' | kubectl apply -f -

dev-apply: dev-secrets ## Apply the k8s manifests and wait for the database
	kubectl apply -k zarf/k8s/dev
	kubectl rollout status statefulset/database -n $(NAMESPACE) --timeout=120s

dev-update: dev-build dev-load dev-apply ## Rebuild the image and (re)deploy it
	kubectl rollout restart deployment/sales -n $(NAMESPACE)

dev-status: ## Show the pods in the cluster
	kubectl get pods -n $(NAMESPACE) -o wide

dev-logs: ## Tail the sales-api logs
	kubectl logs -n $(NAMESPACE) -l app=sales -f --tail=100

# ==============================================================================
# Module support

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	go mod tidy

# ==============================================================================
# Running the API locally

.PHONY: run
run: ## Run the sales API locally (needs a reachable PostgreSQL)
	go run api/sales/main.go

# ==============================================================================
# Code quality check locally

.PHONY: test lint vulncheck fix check

check: lint vulncheck test ## Run all checks

lint: ## Run go vet, go fix and staticcheck
	CGO_ENABLED=0 go vet ./...
	go fix -diff ./...
	go run $(STATICCHECK) -checks=all ./...

vulncheck: ## Scan dependencies for known vulnerabilities
	go run $(GOVULNCHECK) ./...

test: ## Run the tests
	go test -count=1 ./...

fix: ## Apply the modernizations that `make lint` reports
	go fix ./...

# ==============================================================================
