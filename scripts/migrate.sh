#!/bin/sh
set -e

echo "Running migrations"

for file in /srv/migrations/*; do
  echo "Applying $file..."
  psql -h $POSTGRES_DB -U $POSTGRES_USER -d $POSTGRES_DB -f $file
done

echo "Migrations completed"
