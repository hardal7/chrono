package repository

import (
	"context"

	"github.com/hardal7/chrono/internal/config"
	logger "github.com/hardal7/chrono/internal/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getConnString() string {
	return "host=" + config.App.DB_HOST +
		" user=" + config.App.DB_USER +
		" password=" + config.App.DB_PASSWORD +
		" dbname=" + config.App.DB_NAME +
		" port=" + config.App.DB_PORT +
		" sslmode=disable"
}

var DB *pgxpool.Pool

func CreateDBConnection() {
	logger.Info("Connecting to database server at: " + config.App.DB_HOST)
	var err error
	DB, err = pgxpool.New(context.Background(), getConnString())

	if err != nil {
		DB.Close()
		logger.Error("Failed to create connection pool")
		logger.Debug(err.Error())
	}

	if err := DB.Ping(context.Background()); err != nil {
		DB.Close()
		logger.Error("Failed to connect to connection pool")
		logger.Debug(err.Error())
	}

	logger.Info("Connected to database server at: " + config.App.DB_HOST)
}
