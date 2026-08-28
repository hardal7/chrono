package test

import (
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Bootstrap() {
	config.Load()
	logger.Init()
	_, err := db.CreateDBConnection()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
