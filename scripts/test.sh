#!/bin/sh
set -e

echo "Running tests"
make test -C /srv
