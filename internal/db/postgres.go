package db

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB      *pgxpool.Pool
	Queries *db.Queries
)

func CreateDBConnection() (*pgxpool.Pool, error) {
	logger.Info("Connecting to database server", "host", config.App.DBHost)

	cfg, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		return nil, fmt.Errorf("Invalid database connection string: %w", err)
	}

	cfg.ConnConfig.Tracer = queryTracer{}
	DB, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		DB.Close()
		return nil, fmt.Errorf("Failed to create connection pool: %w", err)
	}
	logger.Info("Created connection pool")

	if err := DB.Ping(context.Background()); err != nil {
		DB.Close()
		return nil, fmt.Errorf("Failed to connect to connection pool: %w", err)
	}
	Queries = db.New(DB)

	logger.Info("Connected to database server", "host", config.App.DBHost)
	return DB, nil
}

func getConnectionString() string {
	return "host=" + config.App.DBHost +
		" user=" + config.App.DBUser +
		" password=" + config.App.DBPassword +
		" dbname=" + config.App.DBName +
		" port=" + config.App.DBPort +
		" sslmode=disable"
}

type queryTracer struct{}

func (t queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger.Trace(data.SQL, "requestID", ctx.Value(requestctx.RequestID).(string))
	return ctx
}

func (t queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger.Trace(data.CommandTag.String(), "error", data.Err, "requestID", ctx.Value(requestctx.RequestID).(string))
}
