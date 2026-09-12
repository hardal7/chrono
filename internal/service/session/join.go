package session

import (
	"context"
	"errors"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5"
)

// TODO: Join via link
func Join(ctx context.Context, r dto.JoinSessionRequest) error {
	userID := requestctx.GetUserID(ctx)

	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin new transaction: %w: %w", db.ErrBeginTransaction, err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Warn("rollback transaction")
		}
	}()

	t, err := db.Queries.WithTx(tx).JoinSession(ctx, query.JoinSessionParams{
		UserID:        userID,
		OwnerUsername: r.OwnerUsername,
		SessionName:   r.Name,
	})
	if err != nil {
		return fmt.Errorf("join session: %w: %w", db.ErrRunQuery, err)
	}

	if t.Valid {
		err = db.Queries.WithTx(tx).CreateTopic(ctx, query.CreateTopicParams{
			OwnerID: userID,
			Name:    t.String,
		})
		if err != nil {
			return fmt.Errorf("create topic of session: %w: %w", db.ErrRunQuery, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w: %w", db.ErrCommitTransaction, err)
	}

	return nil
}
