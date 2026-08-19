package topic

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetAll(ctx context.Context) (dto.GetTopicsAllResponse, error) {
	logger.Debug("Getting all topics")
	t, err := conn.Queries.GetTopicsByOwner(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil {
		logger.Debug("Failed to get all topics", err)
		return dto.GetTopicsAllResponse{}, err
	}
	var topics []dto.TopicSelection
	for _, topic := range t {
		topics = append(topics, dto.TopicSelection{Name: topic.Name, TotalTime: int(topic.TotalTimeTrackedSeconds)})
	}
	resp := dto.GetTopicsAllResponse{Topics: topics}
	logger.Debug("Got all topics", "name")
	return resp, nil
}
