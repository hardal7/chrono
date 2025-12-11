package repository

import (
	"context"

	"github.com/hardal7/study/internal/model"
)

func IsDuplicateSession(ctx context.Context, session model.Session) (bool, error) {
	query := "SELECT COUNT(1) FROM sessions WHERE name = $1;"
	var exists int
	err := DB.QueryRow(ctx, query, session.Name).Scan(&exists)

	if exists == 0 {
		return false, err
	} else {
		return true, err
	}
}

func CreateSession(ctx context.Context, session model.Session) error {
	query := "INSERT INTO sessions (name, password, admin_id, user_ids, expiry, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);"
	_, err := DB.Exec(ctx, query, session.Name, session.Password, session.Admin, session.Users, session.Expiry, session.CreatedAt, session.UpdatedAt)

	return err
}
