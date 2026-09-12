package test

import (
	"context"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Bootstrap() {
	config.Load()
	logger.Init()

	_, err := db.CreateDBConnection(context.Background())
	if err != nil {
		logger.Fatal("Failed to connect to DB", err)
	}
	_, err = db.CreateRedisConnection(context.Background())
	if err != nil {
		logger.Fatal("Failed to connect to Redis", err)
	}
}
