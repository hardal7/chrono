package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"uuid"
)

func updateLeaderboard(ctx context.Context) error {
	logger.Info("Updating leaderboard")

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

	snapshotID, err := db.Queries.WithTx(tx).CreateLeaderboardSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create new leaderboard snapshot: %w: %w", db.ErrRunQuery, err)
	}

	users, err := db.Queries.WithTx(tx).GetTopUsers(ctx, query.GetTopUsersParams{
		Limit:     maxRank,
		Cursor:    0,
		MatchName: pgtype.Text{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("Failed to get leaderboard users: %w: %w", db.ErrRunQuery, err)
	}

	var leaderboardUsers []query.CreateLeaderboardUsersParams
	for rank, leaderboardUser := range users {
		leaderboardUsers = append(leaderboardUsers, query.CreateLeaderboardUsersParams{
			SnapshotID: snapshotID,
			UserID:     leaderboardUser.ID,
			Rank:       int32(rank + 1),
		})
	}
	_, err = db.Queries.WithTx(tx).CreateLeaderboardUsers(ctx, leaderboardUsers)
	if err != nil {
		return fmt.Errorf("Failed to create leaderboard users: %w: %w", db.ErrCommitTransaction, err)
	}

	lastUsers, err := db.Queries.WithTx(tx).GetLastLeaderboardUsers(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get last leaderboard users: %w: %w", db.ErrRunQuery, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w: %w", db.ErrCommitTransaction, err)
	}

	topUsers := calculateRanks(users, lastUsers)

	err = updateLeaderboardCache(ctx, topUsers)
	if err != nil {
		return err
	}

	logger.Info("Updated leaderboard")
	return nil
}

func calculateRanks(users []query.User, lastUsers []query.LeaderboardUser) []dto.TopUser {
	lastRanks := make(map[uuid.UUID]int, len(lastUsers))
	for _, lastUser := range lastUsers {
		lastRanks[lastUser.ID] = int(lastUser.Rank)
	}

	var topUsers []dto.TopUser
	for i, user := range users {
		rank := i + 1
		lastRank, ok := lastRanks[user.ID]

		rankChange := rank
		if ok {
			rankChange = lastRank - rank
		}

		topUsers = append(topUsers, dto.TopUser{
			Rank:       rank,
			RankChange: rankChange,
			Username:   user.Username,
			TotalTime:  int(user.TotalTimeTrackedSeconds),
			TodayTime:  int(user.TodayTimeTrackedSeconds),
			AvatarPath: filepath.Join(config.AvatarEndpoint, user.ID.String()),
		})
	}

	return topUsers
}

func updateLeaderboardCache(ctx context.Context, topUsers []dto.TopUser) error {
	body := dto.GetTopUsersResponse{Users: topUsers}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Failed to marshal cache body: %w", err)
	}

	err = db.RDB.Set(ctx, user.TopUsersKey, data, 0).Err()
	if err != nil {
		return fmt.Errorf("Failed to save top users to cache: %w", err)
	}

	return nil
}
