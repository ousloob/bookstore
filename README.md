# bookstore

[![Go](https://img.shields.io/github/go-mod/go-version/ousloob/bookstore)](go.mod)
[![CI](https://github.com/ousloob/bookstore/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ousloob/bookstore/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/ousloob/bookstore)](https://github.com/ousloob/bookstore/releases)
[![License: MIT](https://img.shields.io/github/license/ousloob/bookstore)](LICENSE)

A REST API for an independent bookstore, written in Go. The store manages the books it owns; customers browse the catalog, keep favorites, reserve a book, and buy a physical or digital copy.

## What works today

- HTTP server with graceful shutdown and request logging
- PostgreSQL connection pool
- CI: lint, vulnerability scan, tests and build

## Planned

| Area          | What                                                     |
| ------------- | -------------------------------------------------------- |
| Accounts      | Registration, login, roles                               |
| Catalog       | Manage books, browse with search, sorting and pagination |
| Favorites     | Save books to a list                                     |
| Reservations  | Hold a book for a limited time                           |
| Orders        | Buy physical copies, with safe stock handling            |
| Payments      | Pay for an order, and get refunded                       |
| Digital books | Buy an online edition and download it                    |
| Store tools   | Inventory adjustments, low-stock alerts, sales reports   |

Card numbers are never stored: payment goes through an external provider.

See [ROADMAP.md](ROADMAP.md) for the order of work.

## Tech stack

Go, PostgreSQL, [pgx](https://github.com/jackc/pgx), [caarlos0/env](https://github.com/caarlos0/env) for configuration, [Goose](https://github.com/pressly/goose) for migrations, `net/http` and `log/slog` from the standard library, GitHub Actions.

## Running it

Configuration comes from the environment, with the defaults in the config struct
in [api/sales/main.go](api/sales/main.go). Only the database credentials have no
default, and they are secrets, so they live in [Infisical](https://infisical.com)
rather than in a file:

| Secret        | What                   |
| ------------- | ---------------------- |
| `DB_USER`     | PostgreSQL role        |
| `DB_PASSWORD` | Password for that role |

`infisical run` injects them as environment variables, so nothing in the
application knows Infisical exists. One-time setup:

```sh
infisical login          # stores the credentials in the OS keyring
infisical init           # writes .infisical.json (project id, no secrets)
```

Then, with a PostgreSQL server reachable (local cluster tooling is planned in
the roadmap):

```sh
infisical run -- make run
```

## Development

```sh
make test       # run the tests
make lint       # go vet + go fix -diff + staticcheck
make vulncheck  # govulncheck
make check      # lint + vulncheck + test
make fix        # apply the modernizations that make lint reports
```

These are the checks CI runs, with the tools pinned in the [Makefile](Makefile), so nothing needs installing.

## License

[MIT](LICENSE)
