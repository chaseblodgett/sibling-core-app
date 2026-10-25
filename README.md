# Sibling Core

A family coordination app: task board, expense splitting, and shared documents.

## Stack

- **Backend:** Go, Echo v5, PostgreSQL (via Supabase), [goqu](https://github.com/doug-martin/goqu) for query building
- **Frontend:** React + TypeScript (Vite), MUI
- **Auth:** custom email/password + JWT, owned by the Go backend (not Supabase Auth)
- **Database / Storage:** Supabase
- **Testing:** [testify](https://github.com/stretchr/testify) + [mockery](https://vektra.github.io/mockery/)-generated mocks

## Project layout

```
backend/api/
  cmd/server/          entrypoint (main.go)
  internal/handler/    Echo HTTP handlers, one file per resource
  internal/service/    business logic, interfaces + implementations
  internal/repository/ Postgres data access (goqu query building), interfaces + implementations
  internal/models/     domain structs (Family, Member, Task, Expense, User, ...)
  internal/middleware/ cross-cutting Echo middleware (JWT auth)
  internal/di/         dependency injection wiring
  internal/config/     env-based config loading
  test/mock/           mockery-generated testify mocks, one subfolder per layer
    mock_repository/     mocks of every internal/repository interface
    mock_service/         mocks of every internal/service interface
  .mockery.yml          mockery config — `make mocks` regenerates everything above

frontend/
  src/pages/           route-level pages (login, signup, dashboard, ...)
  src/components/      shared UI components
  src/api/             typed fetch clients for the backend + Supabase client
  src/types/           TypeScript types mirroring backend models
  src/theme/           MUI theme setup

supabase/
  config.toml          local stack config (ports, auth settings, etc.)
  migrations/           SQL migrations — single source of truth for schema

Makefile               local dev commands (install, dev, build, supabase-*)
scripts/dev.sh          runs backend + frontend together, used by `make dev`
```

## Prerequisites

- Go 1.25+
- Node 20+
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (or a compatible runtime —
  OrbStack, Colima) running, for the local Supabase stack
- [Supabase CLI](https://supabase.com/docs/guides/local-development/cli/getting-started)
  (`brew install supabase/tap/supabase`)

## Local Supabase stack (fully offline)

This runs Postgres, Auth, Storage, Realtime, and Studio in Docker — no hosted project or network
access required.

```bash
make supabase-start    # first run pulls Docker images, takes a minute or two
make supabase-status   # prints API URL, DB URL, anon key, service_role key, Studio URL
```

Copy the printed values into your `.env` files (see below) — `DATABASE_URL` and `SUPABASE_URL`
already default to the standard local ports in `.env.example`, so you mainly need the anon and
service_role keys from `make supabase-status`.

Useful extras:

- **Studio** (table editor, Auth users, Storage browser): the URL from `make supabase-status`,
  normally `http://127.0.0.1:54323`
- **Inbucket** (catches emails sent by local Auth instead of delivering them):
  `http://127.0.0.1:54324`
- `make supabase-reset` — drops and recreates the local DB, replaying every migration in
  `supabase/migrations/` from scratch (use this after editing a migration)
- `make supabase-stop` — stops the containers when you're done

Migrations live only in `supabase/migrations/` (CLI-managed, timestamped files). To change the
schema: `supabase migration new <name>`, edit the generated file, then `make supabase-reset` to
apply it locally.

## Hosted Supabase (alternative to Docker)

1. Create a project at [supabase.com](https://supabase.com).
2. **Project Settings → API** — copy:
   - `Project URL` → `SUPABASE_URL` (backend) / `VITE_SUPABASE_URL` (frontend)
   - `anon public` key → `SUPABASE_ANON_KEY` (backend) / `VITE_SUPABASE_ANON_KEY` (frontend)
   - `service_role` key → `SUPABASE_SERVICE_ROLE_KEY` (backend only — **never** put this in the
     frontend `.env`, it bypasses Row Level Security)
3. **Project Settings → Database** — copy the connection string (direct connection, not the
   transaction pooler, since the Go backend holds a persistent connection pool via `pgxpool`) →
   `DATABASE_URL` (backend only).
4. Link the CLI to the project and push the schema: `supabase link --project-ref <ref>` then
   `make db-push`.

Either way, Supabase here is just the Postgres host (and, later, Storage) — user identity lives in
the backend's own `users` table (`supabase/migrations/20261001022021_add_users_auth.sql`), not
Supabase Auth; `members.user_id`, `families.created_by`, etc. all reference it. Storage and RLS
policies aren't wired up yet since the family/task/expense/document service and repository logic
is still stubbed (see Status below).

## Environment files

```bash
cp backend/api/.env.example backend/api/.env
cp frontend/.env.example frontend/.env
```

Fill in the anon/service_role keys from `make supabase-status` (local) or the dashboard (hosted).
Both `.env` files are gitignored. The backend loads `backend/api/.env` automatically on startup
(via `godotenv`); in production, set real environment variables instead and skip the file.

## Running everything locally

```bash
make supabase-start   # once per session — starts the local Supabase stack
make install          # go mod download + npm install
make dev              # runs backend (:8080) and frontend (:5173) together, Ctrl+C stops both
```

`make dev` refuses to start if either `.env` file is missing, so you can't forget to configure
Supabase first. Run `make help` for the full list of targets.

## Testing

```bash
make test     # go test ./... in backend/api
make mocks    # regenerate backend/api/test/mock/** from current interfaces (requires mockery)
```

Backend unit tests use [testify](https://github.com/stretchr/testify) (`assert`/`require`) and
inject mocks for dependency isolation — services get a mocked repository, handlers get a mocked
service. The mocks in `backend/api/test/mock/mock_repository` and `mock_service` are generated by
[mockery](https://vektra.github.io/mockery/) (config: `backend/api/.mockery.yml`) from the
interfaces in `internal/repository` and `internal/service` — **don't hand-edit them**; change the
interface and run `make mocks` instead. Install mockery with `brew install mockery`.

## Status

This is a scaffold: routes, types, and layered structure are in place, but service/repository
logic is stubbed out (`not implemented`) pending feature-by-feature implementation.
