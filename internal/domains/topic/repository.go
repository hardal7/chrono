package topic

import (
	"context"
	"strconv"

	"github.com/hardal7/chrono/internal/db"
)

type Repository interface {
	FindByName(ctx context.Context, name string) (Topic, error)
	FindByID(ctx context.Context, id int) (Topic, error)
	Create(ctx context.Context, topic Topic) error
	Update(ctx context.Context, topic Topic) error
	Delete(ctx context.Context, topic Topic) error
}
type repository struct{}

var Repo Repository = repository{}

const table string = "topics"

func (r repository) FindByName(ctx context.Context, name string) (Topic, error) {
	return db.Get[Topic](ctx, table, "name", name)
}
func (r repository) FindByID(ctx context.Context, id int) (Topic, error) {
	return db.Get[Topic](ctx, table, "id", strconv.Itoa(id))
}
func (r repository) Create(ctx context.Context, topic Topic) error {
	return db.Create(ctx, table, topic)
}
func (r repository) Update(ctx context.Context, topic Topic) error {
	return db.Update(ctx, table, topic)
}
func (r repository) Delete(ctx context.Context, topic Topic) error {
	return db.Delete(ctx, table, topic)
}
