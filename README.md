# Sibling Core

A family coordination app: task board, expense splitting, and shared documents.

## Stack

- **Backend:** Go, Echo v5, PostgreSQL (via Supabase)
- **Frontend:** React + TypeScript (Vite), MUI
- **Database / Auth / Storage:** Supabase

## Project layout

```
backend/api/
  cmd/server/          entrypoint (main.go)
  internal/handler/    Echo HTTP handlers, one file per resource
  internal/service/    business logic, interfaces + implementations
  internal/repository/ Postgres data access, interfaces + implementations
  internal/models/     domain structs (Family, Member, Task, Expense, ...)
  internal/di/         dependency injection wiring
  internal/config/     env-based config loading
  migrations/          SQL migrations

frontend/
  src/pages/           route-level pages (login, signup, dashboard, ...)
  src/components/      shared UI components
  src/api/             typed fetch clients for the backend + Supabase client
  src/types/           TypeScript types mirroring backend models
  src/theme/           MUI theme setup
```

## Prerequisites

- Go 1.25+
- Node 20+
- A Supabase project (or the Supabase CLI running locally)

## Backend setup

```bash
cd backend/api
cp .env.example .env   # fill in DATABASE_URL / SUPABASE_* values
go mod tidy
go run ./cmd/server
```

The server starts on `http://localhost:8080` (see `PORT` in `.env`). Health check: `GET /health`.

### Migrations

SQL migrations live in `backend/api/migrations`. Apply them against your Supabase project with
the Supabase CLI (`supabase db push` / `supabase migration up`) or any Postgres migration tool
that reads plain `.up.sql` / `.down.sql` files (e.g. `golang-migrate`).

## Frontend setup

```bash
cd frontend
cp .env.example .env   # fill in VITE_API_BASE_URL / VITE_SUPABASE_* values
npm install
npm run dev
```

The app starts on `http://localhost:5173` and expects the backend at `VITE_API_BASE_URL`
(defaults to `http://localhost:8080/api/v1`).

## Running both locally

1. Start Supabase (hosted project, or `supabase start` for the local stack).
2. In one terminal: `cd backend/api && go run ./cmd/server`.
3. In another terminal: `cd frontend && npm run dev`.
4. Open `http://localhost:5173`.

## Status

This is a scaffold: routes, types, and layered structure are in place, but service/repository
logic is stubbed out (`not implemented`) pending feature-by-feature implementation.
