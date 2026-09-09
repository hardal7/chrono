#!/bin/sh
set -eu

if [ "${1:-}" = "test" ]; then
    echo "Running migrations in test mode"
else
    echo "Running migrations in normal mode, test migrations skipped"
fi

: "${POSTGRES_HOST:?POSTGRES_HOST is required}"
: "${POSTGRES_PORT:?POSTGRES_PORT is required}"
: "${POSTGRES_USER:?POSTGRES_USER is required}"
: "${POSTGRES_DB:?POSTGRES_DB is required}"
: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"

until pg_isready \
    -h "$POSTGRES_HOST" \
    -p "$POSTGRES_PORT" \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB"
do
    echo "Waiting for PostgreSQL..."
    sleep 1
done

echo "PostgreSQL is ready"

for file in /srv/migrations/*.sql; do
    filename=$(basename "$file")

    if [ "${1:-}" != "test" ] && echo "$filename" | grep -q '^100'; then
        continue
    fi

    echo "Applying $filename..."

    PGPASSWORD="$POSTGRES_PASSWORD" psql \
        -h "$POSTGRES_HOST" \
        -p "$POSTGRES_PORT" \
        -U "$POSTGRES_USER" \
        -d "$POSTGRES_DB" \
        -f "$file"
done

echo "Migrations completed"
