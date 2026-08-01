package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	TrackTime(ctx context.Context, time int, userID string) error
	GetTopUsers(ctx context.Context, cursor int, limit int) ([]User, error)
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
const timeTable string = "total_time_tracked_seconds"

func (r repository) TrackTime(ctx context.Context, time int, userID string) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s + $1 WHERE id = $2;", table, timeTable, timeTable)
	logger.Debug(">", "query", query)
	_, err := db.DB.Exec(ctx, query, time, userID)
	return err
}

func (r repository) GetTopUsers(ctx context.Context, cursor int, limit int) ([]User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s < $1 ORDER_BY %s DESC LIMIT $1;", table, timeTable, timeTable)
	logger.Debug(">", "query", query)
	row, _ := db.DB.Query(ctx, query, cursor, limit)
	users, err := pgx.CollectRows(row, pgx.RowToStructByName[User])
	return users, err
}

func (r repository) FindUserTopic(ctx context.Context, userid int, topicName string) ([]User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s = $2 LIMIT 1;", table, "created_by_userid", "name")
	logger.Debug(">", "query", query)
	row, _ := db.DB.Query(ctx, query, userid, topicName)
	users, err := pgx.CollectRows(row, pgx.RowToStructByName[User])
	return users, err
}

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
