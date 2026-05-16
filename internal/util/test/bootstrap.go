package test

import (
	"github.com/hardal7/chrono/internal/repository"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Bootstrap() {
	logger.Init()
	config.Load()
	repository.CreateDBConnection()
}
