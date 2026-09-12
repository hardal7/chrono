package user

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/location"
	"github.com/hardal7/chrono/internal/service/topic"
	"github.com/hardal7/chrono/internal/util/apierror"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost int = 12

func Register(ctx context.Context, r dto.RegisterUserRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	_, err = db.Queries.GetUserByUsername(ctx, r.Username)
	if !errors.Is(err, pgx.ErrNoRows) {
		if err == nil {
			return fmt.Errorf("user already exists: %w", apierror.ErrAlreadyExists)
		}
		return fmt.Errorf("check if user is duplicate: %w: %w", db.ErrRunQuery, err)
	}

	country := location.IPToCountry(ctx.Value(requestctx.IP).(string))
	userID, err := db.Queries.CreateUser(ctx, query.CreateUserParams{
		Username: r.Username,
		Email:    r.Email,
		Password: string(passwordHash),
		Country: pgtype.Text{
			String: country,
			Valid:  country != "",
		},
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	initAccount(ctx, userID)
	return nil
}

func initAccount(ctx context.Context, userID uuid.UUID) {
	ctx = requestctx.AsUserID(ctx, userID)
	err := topic.InitFirst(ctx)
	if err != nil {
		logger.Err(err).Warn("initialize first topic")
	}
	err = InitAvatar(ctx)
	if err != nil {
		logger.Err(err).Warn("init user avatar")
	}
}
