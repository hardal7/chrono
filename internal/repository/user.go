package repository

import (
	"context"

	"github.com/hardal7/chrono/internal/model"
	"github.com/jackc/pgx/v5"
)

func GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1;"
	row, err := DB.Query(ctx, query, username)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.User])
	return user, err
}
