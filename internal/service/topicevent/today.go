package topicevent

import (
	"context"

	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetToday(ctx context.Context) (dto.GetTopicEventsTodayResponse, error) {
	logger.Info("Getting total time tracked today")
	topicEvents, err := conn.Queries.GetTopicEventsToday(ctx, ctx.Value(middleware.UserID).(int32))
	if err != nil {
		logger.Error("Failed to get topic events today", err)
		return dto.GetTopicEventsTodayResponse{}, err
	}
	var totalTime int
	for i := range topicEvents {
		totalTime += int(topicEvents[i].TimeTrackedSeconds)
	}
	resp := dto.GetTopicEventsTodayResponse{
		TotalTime: totalTime,
	}
	return resp, nil
}
