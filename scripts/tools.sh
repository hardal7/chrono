#!/bin/sh
set -e

echo "Installing tools"
pacman -S docker docker-compose go make # Rather use your systems package manager
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.12.2
