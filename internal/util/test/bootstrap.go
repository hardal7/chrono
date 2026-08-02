package test

import (
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Bootstrap() {
	config.Load()
	logger.Init()
	conn.CreateConnection()
}
