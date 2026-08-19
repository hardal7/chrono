package topicevent

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetToday(ctx context.Context, r dto.GetTopicEventsTodayRequest) (dto.GetTopicEventsTodayResponse, error) {
	logger.Debug("Getting times tracked today")
	userID := ctx.Value(middleware.UserID).(uuid.UUID)
	resp := dto.GetTopicEventsTodayResponse{}

	if len(r.Topics) == 0 {
		events, err := conn.Queries.GetTopicEventsTodayAll(ctx, userID)
		if err != nil {
			logger.Debug("Failed to get topic events today", err)
			return dto.GetTopicEventsTodayResponse{}, err
		}

		for _, event := range events {
			resp.TotalTime += int(event.TimeTrackedSeconds)
		}
		logger.Debug("Got total time tracked today")
		return resp, nil
	}

	for _, topic := range r.Topics {
		events, err := conn.Queries.GetTopicEventsTodayWithTopicName(ctx, db.GetTopicEventsTodayWithTopicNameParams{
			UserID: userID,
			Name:   topic,
		})
		if err != nil {
			logger.Debug("Failed to get topic events today", err)
			return dto.GetTopicEventsTodayResponse{}, err
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

	logger.Debug("Got total time tracked today")
	return resp, nil
}
