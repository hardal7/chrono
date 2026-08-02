#!/bin/sh
set -e

echo "Installing tools"
pacman -S docker docker-compose go make # Rather use your systems package manager
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
sudo curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.12.2
