package topicevent

import (
	"context"
	"time"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/repository"
)

const defaultTopic string = "General"

func Initialize(userID int, ctx context.Context) {
	firstTopic, _ := repository.Find[topic.Topic](ctx, "topics", "name", defaultTopic)
	topicEvent := TopicEvent{
		UserID:      userID,
		TopicID:     firstTopic.ID,
		TimeTracked: 0,
		Date:        int(time.Now().Unix()),
	}
	repository.Create(ctx, topicEvent, "topic_events")
}
