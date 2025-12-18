## Project Structure

The project uses **Domain Driven Design** to enhance maintainability and scalability. The directories are organized according to their features. Below is an overview of the key directories and files:

`api/`: Contains route definitions and API decoders/encoders.\
`middleware/`: Contains custom middleware for request processing.\
`handler/`: Contains business logic for different routes.\
`repository/`: Contains database interaction logic.\
`model/`: Contains Request and Data object definitions.

The general flow of the application is as follows:

1. **API Layer**: Takes a JSON request, decodes it, and creates a Request object (such as `CreateSessionRequest` defined in `model/`).
2. **Middleware Layer**: Processes requests before they reach the handlers.
3. **Handler Layer**: Takes the Request object, processes it according to business logic, and interacts with the Repository layer.
4. **Repository Layer**: Interacts with the database to perform operations, takes or returns Data objects (such as `User` or `Session` defined in `model/`).

Other important files include:

`cmd/main.go`: Entry point of the application.\
`api/request.go`: Decodes incoming API requests into Request objects.\
`repository/connection.go`: Create and manage database connections.\
`config/`: Load application configuration from environment variables or config files.\
`utils/`: Contains utility functions used across the project.

For an overview of the database schema: see `docs/schema.sql`.
