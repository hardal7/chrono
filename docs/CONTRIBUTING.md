## Contributing

Thank you for contributing to this project,
This document will explain all you need to know to start.

## Getting Started

Make sure to have these installed on your system:
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

> [!TIP]
> After reading to end, try implementing your first domain, taking inspiration from the `user` domain. 

## High Level Architecture Design

- Golang with the Chi router is used for the backend
- PostgreSQL is as the database
- SQLc is used to generate code from SQL queries
- Redis is used as the in-memory database
- Traefik is used as the reverse proxy with various other uses
- Victoria Metrics/Logs/vmagent, Vector and Grafana are used as the observability stack
- Docker is used to containerize the services

## Low Level Architecture Design

> [!TIP]
> `domain` refers to a specific area such as `user`, `topic`, `session` etc.\
> `feature` refers to an action on a specific `domain`: `RegisterUser`, `DeleteTopic`

### Overview on the lifecycle of an HTTP request.

HTTP Request -> `api/serve.go` maps to domain routes:

```go
r.Route("/user", UserRoute)
```

`api/domain.go` maps to respective featuer handler:

```go
func UserRoute(r chi.Router) {
	r.Get("/account", GetUserAccountHandler)
	r.Post("/register", RegisterUserHandler)
```

`api/domain.go` Handler runs:

```go
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := user.Register(r.Context(), req)
		processResponse(r.Context(), response{w, resp, err})
	}
}
```

#### `processRequest()`: 
-   Decodes the HTTP JSON body into a struct:

```json
{"email":"mail.com", "username":"user", "password":"secret"}
```

Is decoded into:

```go
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Username string `json:"username" validate:"min=4,max=32"`
	Password string `json:"password" validate:"min=4"`
}
```

- Validates struct against different rules set inside the `dto/feature.go`

`json: "email" validate:"email"` means the JSON body field `email` must be a valid email

- On error returns status

```go
if errors.Is(err, db.ErrNotFound) {
	logger.Debug(err.Error())
	http.Error(w, "Not Found", http.StatusNotFound)
    return
}
```

#### `feature.Service()`
- Calls the respective service of the handler inside `service/domain/feature.go`

Service function has type:

```go
func Register(ctx context.Context, r dto.RegisterUserRequest) (error dto.RegisterUserResponse) {}
```
    
Takes a struct if expecting a request body & returns a struct if expected a response body
    
- UserID of the requester can be retrieved with:

```go
userID := auth.UserID(ctx)
```

- Service interacts with the database using `internal/db` and `internal/db/sqlc` (is `query`):

```go
err := db.Queries.RegisterUser(ctx, query.CreateUserParams{
	userID: userID,
	Name:    r.Username,
})
```

- `Queries` functions are generated from SQL queries inside `sql/queries/domain.sql`:

```sql
-- name: CreateUser :exec
INSERT INTO users(email, username, password)
VALUES($1, $2, $3);
```

- A new domain itself also requires the tables to be created inside `sql/schema.sql`:

```sql
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(64) NOT NULL UNIQUE,
    username VARCHAR(64) NOT NULL UNIQUE,
    password TEXT NOT NULL
```

> Run `make sqlc` to generate the queries after writing them

### `processResponse()`
    
- Handles errors and returns error status if any

- Encodes the struct returned from the service into a JSON body

- Returns the response body and headers

## Other directories

`internal/middleware/`: Functions ran before the HTTP request reaches the handlers\
e.g.: Log request

`internal/runner/`: Functions that are running in the background\
e.g.: Cache frequently requested data
