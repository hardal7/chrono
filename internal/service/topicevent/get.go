package topicevent

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Get(ctx context.Context, r dto.GetTopicEventsRequest) (dto.GetTopicEventsResponse, error) {
	logger.Debug("Getting topic events")
	topicEvents, err := conn.Queries.GetTopicEventsAll(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil {
		logger.Debug("Failed to get topic events", err)
		return dto.GetTopicEventsResponse{}, err
	}

	resp := dto.GetTopicEventsResponse{}
	for i, event := range topicEvents {
		t, _ := conn.Queries.GetTopicByID(ctx, event.TopicID)
		if (r.Topic == "" || r.Topic == t.Name) && (event.CreatedAt.Before(r.FromDate) && event.CreatedAt.After(r.ToDate)) {
			resp.Topics[i] = t.Name
			resp.Dates[i] = event.CreatedAt
			resp.TimesTracked[i] = int(event.TimeTrackedSeconds)
		}
	}

	logger.Debug("Got topic events")
	return resp, nil
}
