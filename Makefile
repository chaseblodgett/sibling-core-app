.DEFAULT_GOAL := help

.PHONY: help install dev backend frontend build test \
	supabase-start supabase-stop supabase-status supabase-reset db-push

help:
	@echo "Targets:"
	@echo "  install          Install backend and frontend dependencies"
	@echo "  dev              Run backend + frontend together (Ctrl+C stops both)"
	@echo "  backend          Run only the Go API server"
	@echo "  frontend         Run only the React dev server"
	@echo "  build            Build backend binary and frontend production bundle"
	@echo "  test             Run backend unit tests (go test ./...)"
	@echo ""
	@echo "  supabase-start   Start the local Supabase stack (Postgres/Auth/Storage/Studio) in Docker"
	@echo "  supabase-stop    Stop the local Supabase stack"
	@echo "  supabase-status  Print local Supabase URLs and API keys"
	@echo "  supabase-reset   Recreate the local DB and replay all migrations"
	@echo "  db-push          Push supabase/migrations to a linked hosted Supabase project"

install:
	cd backend/api && go mod download
	cd frontend && npm install

dev:
	./scripts/dev.sh

backend:
	cd backend/api && go run ./cmd/server

frontend:
	cd frontend && npm run dev

build:
	cd backend/api && go build -o bin/server ./cmd/server
	cd frontend && npm run build

test:
	cd backend/api && go test ./...

supabase-start:
	supabase start

supabase-stop:
	supabase stop

supabase-status:
	supabase status

supabase-reset:
	supabase db reset

db-push:
	supabase db push
