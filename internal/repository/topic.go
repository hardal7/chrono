package repository

import (
	"context"

	"github.com/hardal7/chrono/internal/model"
	"github.com/jackc/pgx/v5"
)

func GetTopicByName(ctx context.Context, name string) (model.Topic, error) {
	query := "SELECT * FROM topics WHERE name = $1 LIMIT 1;"
	row, err := DB.Query(ctx, query, name)
	topic, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.Topic])
	return topic, err
}
