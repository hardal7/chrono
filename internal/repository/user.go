package repository

import (
	"context"

	"github.com/hardal7/study/internal/model"
	"github.com/jackc/pgx/v5"
)

func IsDuplicateUser(ctx context.Context, user model.User) (bool, error) {
	query := "SELECT COUNT(1) FROM users WHERE (email = $1 OR username = $2);"
	var exists int
	err := DB.QueryRow(ctx, query, user.Email, user.Username).Scan(&exists)

	if exists == 0 {
		return false, err
	} else {
		return true, err
	}
}

func CreateUser(ctx context.Context, user model.User) error {
	query := "INSERT INTO users (email, username, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);"
	_, err := DB.Exec(ctx, query, user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt)

	return err
}

func DeleteUser(ctx context.Context, user model.User) error {
	query := "DELETE FROM users WHERE username = $1;"
	_, err := DB.Exec(ctx, query, user.Username)

	return err
}

func UpdateUser(ctx context.Context, user model.User) error {
	query := "UPDATE users SET username = $1, password = $2 WHERE id = $3; UPDATE users SET updated_at = $4"
	_, err := DB.Exec(ctx, query, user.Username, user.Password, user.ID, user.UpdatedAt)

	return err
}

func GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1;"
	row, err := DB.Query(ctx, query, username)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.User])
	return user, err
}

func GetUserByID(ctx context.Context, userid int) (model.User, error) {
	query := "SELECT * FROM users WHERE id = $1 LIMIT 1;"
	row, err := DB.Query(ctx, query, userid)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.User])
	return user, err
}
