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

func GetTopicUserByName(ctx context.Context, topicName string, userID int) (model.TopicUser, error) {
	topicID, _ := GetTopicByName(ctx, topicName)
	query := "SELECT * FROM topics WHERE topic_id = $1 AND user_id = $2 LIMIT 1;"
	row, err := DB.Query(ctx, query, topicID, userID)
	topicUser, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.TopicUser])
	return topicUser, err
}
