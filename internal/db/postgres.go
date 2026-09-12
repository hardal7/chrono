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
	logger.With("host", config.App.PostgresHost).
		Info("Connecting to database server")

	cfg, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		return DB, fmt.Errorf("invalid database connection string: %w", err)
	}

	cfg.ConnConfig.Tracer = queryTracer{}
	DB, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return DB, fmt.Errorf("creating connection pool: %w", err)
	}
	logger.Info("Created connection pool")

	if err := DB.Ping(ctx); err != nil {
		DB.Close()
		return DB, fmt.Errorf("connecting to connection pool: %w", err)
	}
	Queries = db.New(DB)

	logger.With("host", config.App.PostgresHost).
		Info("Connected to database server")
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
	logger.With("requestID", requestctx.GetRequestID(ctx)).
		Trace(data.SQL)
	return ctx
}

func (t queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger.Err(data.Err).
		With("requestID", requestctx.GetRequestID(ctx)).
		Trace(data.CommandTag.String())
}
