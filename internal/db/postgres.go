package db

import (
	"context"

	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB      *pgxpool.Pool
	Queries *db.Queries
)

func CreateDBConnection() {
	logger.Info("Connecting to database server", "host", config.App.DB_HOST)

	cfg, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		logger.Fatal("Invalid database connection string", "error", err)
	}

	cfg.ConnConfig.Tracer = queryTracer{}
	DB, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		DB.Close()
		logger.Fatal("Failed to create connection pool", "error", err)
	}
	logger.Info("Created connection pool")

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

type queryTracer struct{}

func (t queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger.Trace(data.SQL, "requestID", ctx.Value(middleware.RequestID).(string))
	return ctx
}

func (t queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger.Trace(data.CommandTag.String(), "error", data.Err, "requestID", ctx.Value(middleware.RequestID).(string))
}
