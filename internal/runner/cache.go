package runner

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	cachedUsersLimit = 100
	infinity         = 1 << 30
)

func nextMinute() time.Time {
	now := time.Now()
	return now.Truncate(time.Minute).Add(time.Minute)
}

func Cache(ctx context.Context) {
	logger.Info("Started runner", "name", "cache")
	for {
		timer := time.NewTimer(time.Until(nextMinute()))
		<-timer.C
		updateCache(ctx)
	}
}

func updateCache(ctx context.Context) {
	users, err := db.Queries.GetTopUsers(ctx, query.GetTopUsersParams{
		Cursor: infinity,
		Limit:  cachedUsersLimit,
	})
	if err != nil {
		logger.Warn("Failed to get top users")
		return
	}

	topUsers := []dto.TopUser{}
	for i, user := range users {
		topUsers = append(topUsers,
			dto.TopUser{
				Rank:       i + 1,
				RankChange: 1,
				Username:   user.Username,
				TotalTime:  int(user.TotalTimeTrackedSeconds),
				TodayTime:  int(user.TodayTimeTrackedSeconds),
				AvatarPath: filepath.Join(config.AvatarEndpoint, user.ID.String()),
			})
	}

	body := dto.GetTopUsersResponse{Users: topUsers}
	data, err := json.Marshal(body)
	if err != nil {
		logger.Warn("Failed to marshal cache body")
		return
	}

	err = db.RDB.Del(ctx, user.TopUsersKey).Err()
	if err != nil {
		logger.Warn("Failed to delete last top users from cache")
		return
	}

	err = db.RDB.Set(ctx, user.TopUsersKey, data, 0).Err()
	if err != nil {
		logger.Warn("Failed to save top users to cache")
		return
	}
}
