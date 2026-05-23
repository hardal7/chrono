## Project Structure

The project uses **Domain Driven Design** to enhance maintainability and scalability. The directories are organized according to their domains.

Some important files include:

`cmd/main.go`: Entry point of the application.\
`utils/handler/handler.go`: Returns an HTTP handler function that decodes incoming API requests into structs.\

Database interactions are often written standalone but there also are helpers for basic CRUD:
`utils/errors/crud.go`: Custom errors for CRUD
`repository/crud.go`: Generic CRUD database implementation

These are the files that appear commonly in `internal/domains/*`:
`routes.go`: Routing for the handlers
`errors.go`: Custom errors specific to domain business logic
`models.go`: Contains both model and DTO definitions
`repository.go`: Repository definition and implementations for DB operations
Rest of the files are generally handlers for their respective domains.

For an overview of the database schema: see `docs/schema.sql`.
