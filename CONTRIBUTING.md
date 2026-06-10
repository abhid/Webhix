[![RU](https://img.shields.io/badge/lang-ru-blue)](docs/CONTRIBUTING.ru.md)

# Contributing to Webhix

Thanks for your interest. This document covers everything you need to get started.

## Prerequisites

- Go 1.24+
- [sqlc](https://sqlc.dev) — regenerating database queries after schema changes
- [goose](https://github.com/pressly/goose) — running migrations manually
- Node.js 20+ — only if touching the UI

## Project structure

```
cmd/webhix/              entry point (main.go)
internal/
  app/                   application wiring (dependencies, services, HTTP setup)
  cli/                   CLI commands (serve, forward)
  config/                environment-based configuration
  core/                  business logic, repository interfaces
  domain/                domain types and errors
  hub/                   SSE event hub
  server/                HTTP handlers, routing, middleware
  store/                 SQLite implementation (sqlc-generated queries)
    migrations/          goose migration files
    query/               raw SQL queries
    sqlc/                generated Go code — do not edit by hand
  web/                   embedded static assets and UI source
    ui/src/              TypeScript/CSS frontend source
pkg/                     shared utilities
docs/                    documentation
config/                  example env files
```

## Running locally

```sh
# Copy example config
cp config/.env.example config/.env

# Run the server
go run ./cmd/webhix serve

# Run with a custom address
go run ./cmd/webhix serve --addr :9090 --base-url http://localhost:9090
```

## Running all checks

```sh
make ci
```

This runs formatting, linting, vetting, and tests. All checks must pass before opening a PR.

Individual commands:

```sh
make fmt        # format Go code
make lint       # run golangci-lint
make test       # run tests
make web-check  # TypeScript check + ESLint + Prettier
```

## Branches

Follow the naming from [docs/branch-patterns.md](docs/branch-patterns.md).

```
feature/my-feature
fix/some-bug
refactor/cleanup-handler
```

- Branch off `main`
- Keep branches short-lived
- One feature or fix per branch

## Commits

Follow the conventions from [docs/commit-patterns.md](docs/commit-patterns.md).

```
feat(server): add replay endpoint
fix(store): handle null body on insert
refactor(core): extract token generation
```

- English only
- Imperative mood: `add`, `fix`, `remove`
- No period at the end
- No `//nolint` comments — fix the lint issue instead

## Pull requests

- PR title follows the same format as commit messages
- Keep PRs small and focused — one thing per PR
- Add a short description of what changed and why
- All CI checks must pass before merging

## Architecture

The project follows a layered architecture. Dependencies go in one direction only:

```
domain  ←  store
domain  ←  core
domain  ←  server
core    ←  server
```

- `domain` has no dependencies on other internal packages
- `core` defines repository interfaces — it does not import `store`
- `server` depends on `core` interfaces, not on `store` directly
- SQL errors must not leak into `server` — wrap them in `domain` errors at the `store` layer

## Database changes

If you change the schema:

1. Add a migration in `internal/store/migrations/` using goose format
2. Update SQL queries in `internal/store/query/`
3. Regenerate: `sqlc generate`
4. Never edit files in `internal/store/sqlc/` by hand

## UI changes

The frontend is in `internal/web/ui/src/` — vanilla TypeScript, no framework.

```sh
make web-dev    # start Vite dev server
make web-build  # build and embed into binary
make web-check  # lint + type check + prettier
```

## Environment

| Variable            | Default                 | Description                                  |
| ------------------- | ----------------------- | -------------------------------------------- |
| `WEBHIX_BASE_URL`   | `http://localhost:8080` | Public base URL for generated endpoint links |
| `WEBHIX_ADDR`       | `:8080`                 | Listen address                               |
| `WEBHIX_DB_PATH`    | `./data`                | SQLite database directory                    |
| `WEBHIX_PASSWORD`   | —                       | Basic auth password                          |
| `WEBHIX_SECRET_KEY` | —                       | API secret key (Bearer / X-Webhix-Key)       |
