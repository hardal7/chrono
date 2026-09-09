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

func CreateDBConnection(ctx context.Context) (*pgxpool.Pool, error) {
	logger.Info("Connecting to database server", "host", config.App.PostgresHost)

	cfg, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		return DB, fmt.Errorf("Invalid database connection string: %w", err)
	}

	cfg.ConnConfig.Tracer = queryTracer{}
	DB, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return DB, fmt.Errorf("Failed to create connection pool: %w", err)
	}
	logger.Info("Created connection pool")

	if err := DB.Ping(ctx); err != nil {
		DB.Close()
		return DB, fmt.Errorf("Failed to connect to connection pool: %w", err)
	}
	Queries = db.New(DB)

	logger.Info("Connected to database server", "host", config.App.PostgresHost)
	return DB, nil
}

func getConnectionString() string {
	return "host=" + config.App.PostgresHost +
		" port=" + config.App.PostgresPort +
		" user=" + config.App.PostgresUser +
		" dbname=" + config.App.PostgresDB +
		" password=" + config.App.PostgresPassword +
		" sslmode=disable"
}

type queryTracer struct{}

func (t queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger.Trace(data.SQL, "requestID", requestctx.GetRequestID(ctx))
	return ctx
}

func (t queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger.Trace(data.CommandTag.String(), "error", data.Err, "requestID", requestctx.GetRequestID(ctx))
}
