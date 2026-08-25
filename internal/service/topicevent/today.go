package topicevent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func GetToday(ctx context.Context, r dto.GetTopicEventsTodayRequest) (dto.GetTopicEventsTodayResponse, error) {
	userID := ctx.Value(middleware.UserID).(uuid.UUID)
	resp := dto.GetTopicEventsTodayResponse{}

	if len(r.Topics) == 0 {
		events, err := db.Queries.GetTopicEventsTodayAll(ctx, userID)
		if err != nil {
			return resp, fmt.Errorf("Failed to get all topic events today: %w: %w", db.ErrRunQuery, err)
		}

		for _, event := range events {
			resp.TotalTime += int(event.TimeTrackedSeconds)
		}
		return resp, nil
	}

	for _, topic := range r.Topics {
		events, err := db.Queries.GetTopicEventsTodayWithTopicName(ctx, query.GetTopicEventsTodayWithTopicNameParams{
			UserID: userID,
			Name:   topic,
		})
		if err != nil {
			return resp, fmt.Errorf("Failed to get topic events today: %w: %w", db.ErrRunQuery, err)
		}

		var time int
		for _, event := range events {
			time += int(event.TimeTrackedSeconds)
		}
		resp.Topics = append(resp.Topics, dto.TopicEventsToday{
			Name: topic,
			Time: time,
		})
		resp.TotalTime += time
	}

	return resp, nil
}
