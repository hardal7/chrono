package topicevent

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func Get(ctx context.Context, r dto.GetTopicEventsRequest) (dto.GetTopicEventsResponse, error) {
	userID := middleware.UserID(ctx)
	resp := dto.GetTopicEventsResponse{}

	topicEvents, err := db.Queries.GetTopicEventsAll(ctx, userID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get topic events: %w: %w", db.ErrRunQuery, err)
	}

	for i, event := range topicEvents {
		t, _ := db.Queries.GetTopicByID(ctx, event.TopicID)
		if (r.Topic == "" || r.Topic == t.Name) && (event.CreatedAt.Before(r.FromDate) && event.CreatedAt.After(r.ToDate)) {
			resp.Topics[i] = t.Name
			resp.Dates[i] = event.CreatedAt
			resp.TimesTracked[i] = int(event.TimeTrackedSeconds)
		}
	}

	return resp, nil
}
