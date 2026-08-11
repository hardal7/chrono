package main

import (
	"context"

	"github.com/hardal7/chrono/internal/api"
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/runner"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func init() {
	config.Load()
	logger.Init()
}

func main() {
	go runner.NewDay(context.Background())
	conn.CreateConnection()
	api.Serve()
}
