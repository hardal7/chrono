package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/hardal7/chrono/internal/api"
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/runner"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func init() {
	config.Load()
	logger.Init()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.CreateDBConnection(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	} else {
		defer pool.Close()
	}

	rdb, err := db.CreateRedisConnection(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	} else {
		defer rdb.Close()
	}

	go runner.NewDay(ctx)
	go runner.NewWeek(ctx)

	api.Serve(ctx)

	<-ctx.Done()
}
