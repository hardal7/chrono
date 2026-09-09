package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/hardal7/chrono/internal/api"
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/runner"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func init() {
	config.Load()
	logger.Init()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var err error
	var pool *pgxpool.Pool
	for {
		pool, err = db.CreateDBConnection(ctx)

		if err != nil {
			logger.Error(err.Error())
		} else {
			break
		}

		select {
		case <-ctx.Done():
			logger.Info("Shutdown requested:", ctx.Err())
			return
		case <-time.After(time.Second):
		}
	}
	defer pool.Close()

	var rdb *redis.Client
	for {
		rdb, err = db.CreateRedisConnection(ctx)

		if err != nil {
			logger.Error(err.Error())
		} else {
			break
		}

		select {
		case <-ctx.Done():
			logger.Info("Shutdown requested:", ctx.Err())
			return
		case <-time.After(time.Second):
		}
	}
	defer rdb.Close()

	go runner.NewMinute(ctx)
	go runner.NewDay(ctx)
	go runner.NewWeek(ctx)

	api.Serve(ctx)

	<-ctx.Done()
}
