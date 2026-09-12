package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func CreateRedisConnection(ctx context.Context) (*redis.Client, error) {
	logger.With("host", config.App.RedisHost).Info("Connecting to redis server")
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.App.RedisHost + ":" + config.App.RedisPort,
		Password: config.App.RedisPassword,
		DB:       0,
	})

	err := RDB.Ping(ctx).Err()
	if err != nil {
		err = errors.Join(err, RDB.Close())
		return RDB, fmt.Errorf("ping server: %w", err)
	}

	logger.With("host", config.App.RedisHost).Info("Connected to redis server")
	return RDB, nil
}
