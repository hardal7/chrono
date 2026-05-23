package user

import (
	"context"
	"strconv"

	"github.com/hardal7/chrono/internal/db"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int) (User, error)
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, user User) error
}
type repository struct{}

var Repo Repository = repository{}

const table string = "users"

func (r repository) FindByUsername(ctx context.Context, username string) (User, error) {
	return db.Get[User](ctx, table, "username", username)
}
func (r repository) FindByEmail(ctx context.Context, email string) (User, error) {
	return db.Get[User](ctx, table, "email", email)
}
func (r repository) FindByID(ctx context.Context, id int) (User, error) {
	return db.Get[User](ctx, table, "id", strconv.Itoa(id))
}
func (r repository) Create(ctx context.Context, user User) error {
	return db.Create(ctx, table, user)
}
func (r repository) Update(ctx context.Context, user User) error {
	return db.Update(ctx, table, user)
}
func (r repository) Delete(ctx context.Context, user User) error {
	return db.Delete(ctx, table, user)
}
