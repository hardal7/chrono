package topicevent

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5"
)

func Track(ctx context.Context, r dto.TrackTopicEventRequest) error {
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

	topic, err := db.Queries.WithTx(tx).GetTopicByOwnerAndName(ctx, query.GetTopicByOwnerAndNameParams{
		Name:    r.Topic,
		OwnerID: userID,
	})
	if err != nil {
		return fmt.Errorf("get topic by username: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.WithTx(tx).TrackTopicTime(ctx, query.TrackTopicTimeParams{
		ID:          topic.ID,
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		return fmt.Errorf("track topic time: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.WithTx(tx).TrackUserTime(ctx, query.TrackUserTimeParams{
		ID:          userID,
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		return fmt.Errorf("track user time: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.WithTx(tx).CreateTopicEvent(ctx, query.CreateTopicEventParams{
		UserID:             userID,
		TopicID:            topic.ID,
		TimeTrackedSeconds: int32(r.TimeSeconds),
		CreatedAt:          r.Date,
	})
	if err != nil {
		return fmt.Errorf("create topic event: %w: %w", db.ErrRunQuery, err)
	}

	s, err := db.Queries.WithTx(tx).GetJoinedSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("get joined sessions: %w: %w", db.ErrRunQuery, err)
	}

	for _, session := range s {
		if session.Topic.Valid && (r.Topic != session.Topic.String) {
			continue
		}

		err = trackSessionTime(ctx, tx, userID, session, r)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w: %w", db.ErrCommitTransaction, err)
	}

	return nil
}

func trackSessionTime(ctx context.Context, tx pgx.Tx, userID uuid.UUID, session query.Session, r dto.TrackTopicEventRequest) error {
	err := db.Queries.WithTx(tx).TrackSessionTime(ctx, query.TrackSessionTimeParams{
		ID:          session.ID,
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		return fmt.Errorf("track session time: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.WithTx(tx).TrackSessionParticipantTime(ctx, query.TrackSessionParticipantTimeParams{
		UserID:      userID,
		SessionID:   session.ID,
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		return fmt.Errorf("track session participant time: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
