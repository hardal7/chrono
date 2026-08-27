package db

import (
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func CreateRedisConnection() {
	logger.Info("Connecting to redis server", "host", config.App.REDIS_HOST)
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.App.REDIS_HOST + ":" + config.App.REDIS_PORT,
		Password: config.App.REDIS_PASSWORD,
		DB:       0,
	})

	logger.Info("Connected to redis server", "host", config.App.REDIS_HOST)
}
