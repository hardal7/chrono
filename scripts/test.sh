#!/bin/sh
set -e

echo "Initializing test container"
go install github.com/joho/godotenv/cmd/godotenv@latest

echo "Running tests"
make test -C /srv

echo "Initialization completed"
