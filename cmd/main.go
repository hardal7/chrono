package main

import (
	"github.com/hardal7/chrono/internal/api"
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func init() {
	logger.Init()
	config.Load()
}

func main() {
	db.CreateConnection()
	api.Serve()
}
