package user

import (
	"context"
	"time"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/util/logger"
)

const firstTopic string = "General"

func InitTopic(u User) {
	t := topic.Topic{
		Name:      firstTopic,
		CreatedBy: u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := topic.Repo.Create(context.Background(), t); err != nil {
		logger.Warn("Failed to initialize first topic", "username", u.Username)
		return
	} else {
		logger.Info("Initialized first topic", "username", u.Username)
	}
}
