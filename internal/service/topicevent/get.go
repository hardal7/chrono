package topicevent

import (
	"context"

	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Get(ctx context.Context, r dto.GetTopicEventsRequest) (dto.GetTopicEventsResponse, error) {
	logger.Info("Getting topic events")
	topicEvents, err := conn.Queries.GetAllTopicEvents(ctx, ctx.Value(middleware.UserID).(int32))
	if err != nil {
		logger.Error("Failed to get topic events", err)
		return dto.GetTopicEventsResponse{}, err
	}
	resp := dto.GetTopicEventsResponse{}
	for i, v := range topicEvents {
		t, _ := conn.Queries.GetTopicByID(ctx, v.TopicID)
		// TODO: Also check the dates
		if t.Name != "" && t.Name == r.Topic {
			resp.Topics[i] = t.Name
			resp.Dates[i] = v.Date.Time
			resp.TimesTracked[i] = int(v.TimeTrackedSeconds)
		}
	}
	logger.Info("Got topic events")
	return resp, nil
}
