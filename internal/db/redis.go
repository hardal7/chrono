package db

import (
	"context"
	"errors"

	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func CreateRedisConnection(ctx context.Context) (*redis.Client, error) {
	logger.Info("Connecting to redis server", "host", config.App.RedisHost)
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.App.RedisHost + ":" + config.App.RedisPort,
		Password: config.App.RedisPassword,
		DB:       0,
	})

	err := RDB.Ping(ctx).Err()
	if err != nil {
		RDB.Close()
		return RDB, errors.New("Failed to ping redis server")
	}

	logger.Info("Connected to redis server", "host", config.App.RedisHost)
	return RDB, nil
}
