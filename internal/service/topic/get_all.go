package topic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func GetAll(ctx context.Context) (dto.GetTopicsAllResponse, error) {
	resp := dto.GetTopicsAllResponse{}

	t, err := db.Queries.GetTopicsByOwner(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil {
		return resp, fmt.Errorf("Failed to get all topics: %w: %w", db.ErrRunQuery, err)
	}

	var topics []dto.TopicSelection
	for _, topic := range t {
		topics = append(topics, dto.TopicSelection{Name: topic.Name, TotalTime: int(topic.TotalTimeTrackedSeconds)})
	}
	resp = dto.GetTopicsAllResponse{Topics: topics}

	return resp, nil
}
