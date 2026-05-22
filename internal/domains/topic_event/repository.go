package topicevent

import (
	"context"
	"strconv"

	"github.com/hardal7/chrono/internal/db"
)

type Repository interface {
	GetAll(ctx context.Context, id int) ([]TopicEvent, error)
	Create(ctx context.Context, topicevent TopicEvent) error
	Update(ctx context.Context, topicevent TopicEvent) error
	Delete(ctx context.Context, topicevent TopicEvent) error
}
type repository struct{}

var Repo Repository = repository{}

const table string = "topic_events"

func (r repository) GetAll(ctx context.Context, id int) ([]TopicEvent, error) {
	return db.GetMultiple[TopicEvent](ctx, table, "id", strconv.Itoa(id))
}
func (r repository) Create(ctx context.Context, topicevent TopicEvent) error {
	return db.Create(ctx, table, topicevent)
}
func (r repository) Update(ctx context.Context, topicevent TopicEvent) error {
	return db.Update(ctx, table, topicevent)
}
func (r repository) Delete(ctx context.Context, topicevent TopicEvent) error {
	return db.Delete(ctx, table, topicevent)
}
