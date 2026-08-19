package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/service/location"
	"github.com/hardal7/chrono/internal/service/topic"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost int = 12

func Register(ctx context.Context, r dto.RegisterUserRequest) error {
	logger.Debug("Registering user", "username", r.Username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcryptCost)
	if err != nil {
		logger.Warn("Failed to hash password", err)
		return err
	}
	_, err = conn.Queries.GetUserByUsername(ctx, r.Username)
	if err != pgx.ErrNoRows {
		if err == nil {
			logger.Debug("User already exists")
			// TODO: Custom error types
			return errors.New("user already exists")
		} else {
			logger.Warn("Failed to check if user is duplicate", err)
			return err
		}
	}

	country := location.IPToCountry(ctx.Value(middleware.IP).(string))
	err = conn.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: r.Username,
		Email:    r.Email,
		Password: string(passwordHash),
		Country: pgtype.Text{
			String: country,
			Valid:  country != "",
		},
	})
	if err != nil {
		logger.Debug("Failed to create user", err)
		return err
	}
	logger.Debug("Registered user", "username", r.Username)

	u, err := conn.Queries.GetUserByUsername(ctx, r.Username)
	if err != nil {
		logger.Warn("Failed to get created user", "username", r.Username)
	}
	initAccount(ctx, u.ID)
	return nil
}

func initAccount(ctx context.Context, userID uuid.UUID) {
	ctx = context.WithValue(ctx, middleware.UserID, userID)
	err := topic.InitFirst(ctx)
	if err != nil {
		logger.Warn("Failed to initialize first topic")
	}
	err = InitAvatar(ctx)
	if err != nil {
		logger.Warn("Failed to init user avatar", err)
	}
}
