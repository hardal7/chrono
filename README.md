## Chrono Backend

This is the backend service for the [Chrono mobile application](https://github.com/hardal7/chrono-app). The infastructure is built using Go and Postgres.

## Running

### Requirements
- [Docker](https://docs.docker.com/get-started/get-docker/)

Move the example environment file to `.env`:

    mv .env.example .env

Run the volumes script in `scripts` to create the volumes in the host machine:

    ./scripts/volumes.sh

Then run the backend service using Docker Compose, use the following command (run with root privileges):

    make dev-up

## Contributing

All contributions are welcome, kindly open a new issue or pull request. To get started contributing, read the [CONTRIBUTING.md](https://github.com/hardal7/chrono/blob/main/docs/CONTRIBUTING.md) file for on overview on the overall project structure.
