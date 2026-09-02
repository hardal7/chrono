package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	nextWeek = 1
	maxRank  = 100
)

func NewWeek(ctx context.Context) {
	logger.Info("Started runner", "name", "new_week")
	for {
		timer := time.NewTimer(time.Until(retrieveDate(nextWeek)))
		<-timer.C
		err := updateLeaderboard(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
		err = resetWeekTimes(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func updateLeaderboard(ctx context.Context) error {
	logger.Info("Updating leaderboard")

	users, err := db.Queries.GetTopUsers(ctx, query.GetTopUsersParams{
		Limit:     maxRank,
		Cursor:    0,
		MatchName: pgtype.Text{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("Failed to get leaderboard users: %w: %w", db.ErrRunQuery, err)
	}

	lastUsers, err := db.Queries.GetLastLeaderboardUsers(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get last leaderboard users: %w: %w", db.ErrRunQuery, err)
	}

	snapshot, err := db.Queries.CreateLeaderboardSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create new leaderboard snapshot: %w: %w", db.ErrRunQuery, err)
	}

	lastRanks := make(map[uuid.UUID]int, len(lastUsers))
	for _, lastUser := range lastUsers {
		lastRanks[lastUser.ID] = int(lastUser.Rank)
	}

	var leaderboardUsers []query.LeaderboardUser
	for i, user := range users {
		rank := i + 1
		lastRank, ok := lastRanks[user.ID]

		rankChange := rank
		if ok {
			rankChange = lastRank - rank
		}

		leaderboardUsers = append(leaderboardUsers, query.LeaderboardUser{
			SnapshotID: snapshot,
			UserID:     user.ID,
			Rank:       int32(rank),
			RankChange: int32(rankChange),
		})
	}

	logger.Info("Updated leaderboard")
	return nil
}

func resetWeekTimes(ctx context.Context) error {
	logger.Info("Reseting times tracked this week")

	err := db.Queries.ResetUserTimeTrackedWeek(ctx)
	if err != nil {
		return fmt.Errorf("Failed to reset time tracked for this week: %w: %w", db.ErrRunQuery, err)
	}

	logger.Info("Reset times tracked this week")
	return nil
}
