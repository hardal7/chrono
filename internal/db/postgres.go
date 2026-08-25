package db

import (
	"context"

	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB      *pgxpool.Pool
	Queries *db.Queries
)

func CreateDBConnection() {
	logger.Info("Connecting to database server", "host", config.App.DB_HOST)
	var err error
	DB, err = pgxpool.New(context.Background(), getConnectionString())
	if err != nil {
		DB.Close()
		logger.Fatal("Failed to create connection pool", "error", err)
	}
	logger.Info("Created connection pool")
	logger.Info("Connecting to database server")
	if err := DB.Ping(context.Background()); err != nil {
		DB.Close()
		logger.Fatal("Failed to connect to connection pool", "error", err)
	}
	Queries = db.New(DB)
	logger.Info("Connected to database server", "host", config.App.DB_HOST)
}

func getConnectionString() string {
	return "host=" + config.App.DB_HOST +
		" user=" + config.App.DB_USER +
		" password=" + config.App.DB_PASSWORD +
		" dbname=" + config.App.DB_NAME +
		" port=" + config.App.DB_PORT +
		" sslmode=disable"
}
