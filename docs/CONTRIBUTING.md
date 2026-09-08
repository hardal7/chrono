## Contributing

Thank you for contributing to this project,
This document will explain all you need to know to start.

## Getting Started

Make sure to have these tools installed on your system:
- [Go Language](https://go.dev/dl/)
- [Makefile](https://community.chocolatey.org/packages/make)
- [Docker](https://docs.docker.com/get-started/get-docker/)

Move the example environment file to `.env`:
```sh
mv .env.example .env
```

## Running the application
Running the development environment:
```sh
make dev-up
```
Running the tests:
```sh
make test-up
```
If deploying to production:
```sh
make prod-up
```

## High Level Architecture Design

- Golang with the Chi router is used for the backend
- PostgreSQL is used as the database
- SQLc is used to generate code from SQL queries
- Redis is used as the in-memory database
- Traefik is used as the reverse proxy with various other uses
- Victoria Metrics/Logs/vmagent, Vector and Grafana are used as the observability stack
- Docker is used to containerize the services
