package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

func NewMinute(ctx context.Context) {
	logger.Info("Started runner", "name", "new_minute")
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
	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("Failed to begin new transaction: %w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Warn("Failed to rollback transaction")
		}
	}()

	err = db.Queries.WithTx(tx).DeleteExpiredSessions(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete expired sessions: %w: %w", db.ErrRunQuery, err)
	}
	err = db.Queries.WithTx(tx).DeleteExpiredSessionParticipants(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete expired session participants: %w: %w", db.ErrRunQuery, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w: %w", db.ErrCommitTransaction, err)
	}
	return nil
}
