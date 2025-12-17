package repository

import (
	"context"

	"github.com/hardal7/chrono/internal/model"
	"github.com/jackc/pgx/v5"
)

func IsDuplicateTopic(ctx context.Context, topic model.Topic) (bool, error) {
	query := "SELECT COUNT(1) FROM topics WHERE name = $1;"
	var exists int
	err := DB.QueryRow(ctx, query, topic.Name).Scan(&exists)

	if exists == 0 {
		return false, err
	} else {
		return true, err
	}
}

func CreateTopic(ctx context.Context, topic model.Topic) error {
	query := "INSERT INTO topics (name, created_at, updated_at) VALUES ($1, $2, $3);"
	_, err := DB.Exec(ctx, query, topic.Name, topic.CreatedAt, topic.UpdatedAt)

	return err
}

func DeleteTopic(ctx context.Context, topic model.Topic) error {
	query := "DELETE FROM topics WHERE name = $1;"
	_, err := DB.Exec(ctx, query, topic.Name)

	return err
}

func UpdateTopic(ctx context.Context, topic model.Topic) error {
	query := "UPDATE topics SET name = $1 WHERE id = $2;"
	_, err := DB.Exec(ctx, query, topic.Name, topic.ID)

	return err
}

func GetTopicByName(ctx context.Context, name string) (model.Topic, error) {
	query := "SELECT * FROM topics WHERE name = $1 LIMIT 1;"
	row, err := DB.Query(ctx, query, name)
	topic, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.Topic])
	return topic, err
}
