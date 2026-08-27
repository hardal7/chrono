package user

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func EditAccount(ctx context.Context, r dto.EditUserAccountRequest) error {
	userID := middleware.UserID(ctx)

	var password string
	if r.NewPassword != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcryptCost)
		if err != nil {
			return fmt.Errorf("Failed to hash password: %w", err)
		}
		password = string(passwordHash)
	}

	err := db.Queries.UpdateUser(ctx, query.UpdateUserParams{
		ID:       userID,
		Username: r.NewUsername,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("Failed to update user: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}

func DeleteAccount(ctx context.Context) error {
	userID := middleware.UserID(ctx)

	err := db.Queries.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("Failed to delete user: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}

func GetAccount(ctx context.Context) (dto.GetUserAccountResponse, error) {
	userID := middleware.UserID(ctx)
	resp := dto.GetUserAccountResponse{}

	u, err := db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get user: %w: %w", db.ErrRunQuery, err)
	}

	resp = dto.GetUserAccountResponse{
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Time,
	}

	return resp, nil
}
