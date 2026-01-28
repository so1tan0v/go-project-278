#!/bin/sh
set -eu

echo "[run.sh] Starting service"

echo "[run.sh] Running DB migrations"
goose -dir ./db/migrations postgres "${DATABASE_URL:?DATABASE_URL is required}" up

echo "[run.sh] Starting Go app"
exec /app/bin/app