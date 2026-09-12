package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5"
)

func NewMinute(ctx context.Context) {
	logger.With("name", "new_minute").Info("Started runner")
	for {
		timer := time.NewTimer(time.Minute)
		<-timer.C
		err := cleanExpiredSessions(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func cleanExpiredSessions(ctx context.Context) error {
	requestID := "SESSION_CLEANUP"
	ctx = context.WithValue(ctx, requestctx.RequestID, requestID)

	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin new transaction: %w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Err(err).Warn("rollback transaction")
		}
	}()

	err = db.Queries.WithTx(tx).DeleteExpiredSessions(ctx)
	if err != nil {
		return fmt.Errorf("delete expired sessions: %w: %w", db.ErrRunQuery, err)
	}
	err = db.Queries.WithTx(tx).DeleteExpiredSessionParticipants(ctx)
	if err != nil {
		return fmt.Errorf("delete expired session participants: %w: %w", db.ErrRunQuery, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w: %w", db.ErrCommitTransaction, err)
	}
	return nil
}
