#!/usr/bin/env bash
# Runs the Go backend and the React frontend together for local development.
# Ctrl+C stops both.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ ! -f "$ROOT_DIR/backend/api/.env" ]; then
  echo "Missing backend/api/.env — copy backend/api/.env.example and fill in Supabase values first." >&2
  exit 1
fi

if [ ! -f "$ROOT_DIR/frontend/.env" ]; then
  echo "Missing frontend/.env — copy frontend/.env.example and fill in Supabase values first." >&2
  exit 1
fi

cleanup() {
  echo "" && echo "Stopping dev servers..."
  kill 0
}
trap cleanup EXIT INT TERM

(cd "$ROOT_DIR/backend/api" && go run ./cmd/server) &
(cd "$ROOT_DIR/frontend" && npm run dev) &

wait
