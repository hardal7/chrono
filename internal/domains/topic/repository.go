package topic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	TrackTime(ctx context.Context, time int, topicID string) error
	FindUserTopic(ctx context.Context, userid int, topicName string) (Topic, error)
	FindByName(ctx context.Context, name string) (Topic, error)
	FindByID(ctx context.Context, id int) (Topic, error)
	Create(ctx context.Context, topic Topic) error
	Update(ctx context.Context, topic Topic) error
	Delete(ctx context.Context, topic Topic) error
}
type repository struct{}

var Repo Repository = repository{}

const table string = "topics"

func (r repository) TrackTime(ctx context.Context, time int, topicID string) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s + $1 WHERE id = $2;", table, "total_time_tracked_seconds", "total_time_tracked_seconds")
	logger.Debug(">", "query", query)
	_, err := db.DB.Exec(ctx, query, time, topicID)
	return err
}

func (r repository) FindUserTopic(ctx context.Context, userid int, topicName string) (Topic, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s = $2 LIMIT 1;", table, "created_by_userid", "name")
	logger.Debug(">", "query", query)
	row, _ := db.DB.Query(ctx, query, userid, topicName)
	topic, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Topic])
	return topic, err
}

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
