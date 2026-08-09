package user

import (
	"context"
	"errors"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost int = 12

func Register(ctx context.Context, r dto.RegisterUserRequest) error {
	logger.Info("Registering user", "username", r.Username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcryptCost)
	if err != nil {
		logger.Error("Failed to hash password", err)
		return err
	}
	_, err = conn.Queries.GetUserByUsername(ctx, r.Username)
	if err != pgx.ErrNoRows {
		if err == nil {
			logger.Error("User already exists")
			// TODO: Custom error types
			return errors.New("user already exists")
		} else {
			logger.Error("Failed to check if user is duplicate", err)
			return err
		}
	}
	err = conn.Queries.CreateUser(ctx, db.CreateUserParams{
		Username:           r.Username,
		Email:              r.Email,
		Password:           string(passwordHash),
		TimeTrackedSeconds: 0,
	})
	if err != nil {
		logger.Error("Failed to create user", err)
		return err
	}
	logger.Info("Registered user", "username", r.Username)
	return nil
}
