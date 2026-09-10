#!/bin/sh
set -e

echo "Installing tools"
pacman -S docker docker-compose kubectl minikube helm go make # Rather use your systems package manager
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install mvdan.cc/gofumpt@latest
sudo curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin
