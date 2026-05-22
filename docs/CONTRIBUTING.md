## Project Structure

The project uses **Domain Driven Design** to enhance maintainability and scalability. The directories are organized according to their domains.

Some important files include:

`cmd/main.go`: Entry point of the application.\
`utils/handler/handler.go`: Returns an HTTP handler function that decodes incoming API requests into structs.\
`repository/connection.go`: Establishes and manages database connections.\

Implementations are often standalone but there are helpers for basic CRUD:
`utils/errors/crud.go`: Custom errors for CRUD
`repository/crud.go`: Generic CRUD database implementation
Take a look into the code for one of the handlers inside `internal/domains` to get an idea.

For an overview of the database schema: see `docs/schema.sql`.
