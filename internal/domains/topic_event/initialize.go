package topicevent

import (
	"context"
	"time"

	"github.com/hardal7/chrono/internal/domains/topic"
)

const defaultTopic string = "General"

func Initialize(userID int, ctx context.Context) error {
	firstTopic, _ := topic.Repo.FindByName(ctx, defaultTopic)
	topicEvent := TopicEvent{
		UserID:      userID,
		TopicID:     firstTopic.ID,
		TimeTracked: 0,
		Date:        int(time.Now().Unix()),
	}
	err := Repo.Create(ctx, topicEvent)
	return err
}
