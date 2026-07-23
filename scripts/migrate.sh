#!/bin/sh
set -e

if [ "$1" = "test" ]; then
    echo "Running migrations in test mode"
else
    echo "Running migrations in normal mode, test migrations skipped"
fi


for file in /srv/migrations/*.sql; do
    filename=$(basename "$file")

    if [ "$1" != "test" ] && echo "$filename" | grep -q '^100'; then
        continue
    fi

    echo "Applying $filename..."
    psql -h $POSTGRES_DB -U $POSTGRES_USER -d $POSTGRES_DB -f $file
done

echo "Migrations completed"
