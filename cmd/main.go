package main

import (
	"github.com/hardal7/chrono/internal/api"
	"github.com/hardal7/chrono/internal/config"
	"github.com/hardal7/chrono/internal/repository"
	logger "github.com/hardal7/chrono/internal/util"
)

func init() {
	logger.Init()
	config.Load()
}

func main() {
	repository.CreateDBConnection()
	api.RunAPIServer()
}
