#!/bin/sh
# Applies our migrations, then River's, as `make migrate` does.
set -eu
goose -dir /app/db/migrations postgres "$DATABASE_URL" up
river migrate-up --line main --database-url "$DATABASE_URL"
