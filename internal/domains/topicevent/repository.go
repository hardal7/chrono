package topicevent

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetToday(ctx context.Context, userid int) ([]TopicEvent, error)
	GetAll(ctx context.Context, id int) ([]TopicEvent, error)
	Create(ctx context.Context, topicevent TopicEvent) error
	Update(ctx context.Context, topicevent TopicEvent) error
	Delete(ctx context.Context, topicevent TopicEvent) error
}
type repository struct{}

var Repo Repository = repository{}

const table string = "topic_events"

func (r repository) GetToday(ctx context.Context, userid int) ([]TopicEvent, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND DATE(date) = CURRENT_DATE;", table, "user_id")
	logger.Debug(">", "query", query)
	row, _ := db.DB.Query(ctx, query, userid)
	topicEvents, err := pgx.CollectRows(row, pgx.RowToStructByName[TopicEvent])
	return topicEvents, err
}
func (r repository) GetAll(ctx context.Context, userid int) ([]TopicEvent, error) {
	return db.GetMultiple[TopicEvent](ctx, table, "userid", strconv.Itoa(userid))
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
