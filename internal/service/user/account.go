package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func EditAccount(ctx context.Context, r dto.EditUserAccountRequest) error {
	u, err := conn.Queries.GetUserByID(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil {
		logger.Error("Failed to get user", err)
		return err
	}
	logger.Info("Editing account details", "username", u.Username)
	if r.DeleteAccount {
		logger.Info("Deleting account", "username", u.Username)
		err := conn.Queries.DeleteUser(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
		if err != nil {
			logger.Error("Failed to delete user", err)
			return err
		}
		logger.Info("Deleted account", "username", u.Username)
		return nil
	}
	if r.NewUsername != "" {
		u.Username = r.NewUsername
		logger.Info("Changing username", "username", u.Username, "newUsername", r.NewUsername)
		_, err := conn.Queries.GetUserByUsername(ctx, r.NewUsername)
		if err != pgx.ErrNoRows {
			logger.Error("Account with username exists")
			return errors.New("account with username exists")
		}
	}
	if r.NewPassword != "" {
		logger.Info("Changing account password")
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcryptCost)
		if err != nil {
			logger.Error("Failed to hash password", err)
			return err
		}
		u.Password = string(passwordHash)
	}
	err = conn.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		logger.Error("Failed to update user", err)
		return err
	}
	logger.Info("Edited account details")
	return nil
}

func GetAccount(ctx context.Context) (dto.GetUserAccountResponse, error) {
	u, err := conn.Queries.GetUserByID(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil {
		logger.Error("Failed to get user", err)
		return dto.GetUserAccountResponse{}, err
	}
	resp := dto.GetUserAccountResponse{
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Time,
	}

	return resp, nil
}
