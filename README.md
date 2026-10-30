# bookstore

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/ousloob/bookstore/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/ousloob/bookstore/tree/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/ousloob/bookstore)](https://goreportcard.com/report/github.com/ousloob/bookstore)
[![License: MIT](https://img.shields.io/github/license/ousloob/bookstore)](LICENSE)

A REST API for an independent bookstore, written in Go. The store manages the books it owns; customers can browse the catalogue, keep favorites, reserve a book, and buy a physical or digital copy.

> **Status: early development.** The project skeleton is in place, but most features below are still planned. See [ROADMAP.md](ROADMAP.md) for what is next.

## Who it is for

- **Bookstore owners and small shops** who want a simple API for their catalogue, stock and reservations, instead of a heavy e-commerce platform.
- **Go developers** who want a readable example of a layered Go service (handlers, business logic, storage) with PostgreSQL.

## What works today

- HTTP server with graceful shutdown, configurable timeouts and structured logging (`log/slog`)
- Request logging middleware
- Configuration through environment variables and flags
- PostgreSQL connection support (`sqlx`, `lib/pq`)
- CI with tests, `go vet`, `staticcheck` and `govulncheck`

## In progress

- User creation in the business layer, with password hashing (bcrypt). It is not exposed through the API yet, and the database schema is still to be written.

## Planned

| Area          | What                                                             |
| ------------- | ---------------------------------------------------------------- |
| Catalogue     | Add, edit, list, search and paginate books (store staff)         |
| Accounts      | Customer registration and login, roles (customer, staff)         |
| Favorites     | Customers save books to a favorites list                         |
| Reservations  | Hold a book for a limited time, then release it automatically    |
| Orders        | Buy physical copies, with safe stock handling                    |
| Digital books | Buy an online edition and download it through time-limited links |
| Store tools   | Inventory adjustments, low-stock alerts, sales reports           |

Payments are **out of scope** for now: orders are recorded, and card payments would be handed to an external provider later. This API will never store card numbers.

## Tech stack

- **Language:** Go 1.27+
- **HTTP:** [chi](https://github.com/go-chi/chi)
- **Database:** PostgreSQL, accessed with `sqlx`
- **Passwords:** bcrypt
- **Configuration:** [ardanlabs/conf](https://github.com/ardanlabs/conf)
- **CI:** CircleCI

## Getting started

Requirements: Go 1.27 or later. A PostgreSQL server is needed for anything that touches the database (the database schema and migrations are still to be added).

```sh
git clone https://github.com/ousloob/bookstore.git
cd bookstore
make run
```

The server listens on port 3000 by default:

```sh
curl http://localhost:3000/
```

Configuration (web timeouts, listen address, database host, name and credentials) is defined in [`api/sales/main.go`](api/sales/main.go).

## Development

```sh
make test       # run the tests
make lint       # go vet + staticcheck
make vulncheck  # govulncheck
make check      # all of the above
```

## Project layout

```
api/sales/              entry point, routing and handlers
business/domain/        business logic and storage, one package per domain (books, users)
business/sys/           shared support code (database, validation)
support/web/mid/        HTTP middleware
```

## Data

This project uses **made-up sample data only**. No real customer or personal data belongs in this repository.

## Contributing

Issues and pull requests are welcome. Please open an issue first for larger changes.

## License

[MIT](LICENSE)
