package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/redis/go-redis/v9"
)

const otpExpiration = time.Minute * 5

func GenerateOTP(ctx context.Context) (string, error) {
	token, err := generateToken()
	if err != nil {
		return token, fmt.Errorf("generate OTP token: %w", err)
	}

	err = db.RDB.Set(ctx, hashToken(token), requestctx.GetUserID(ctx).String(), otpExpiration).Err()
	if err != nil {
		return token, fmt.Errorf("set OTP token values: %w", err)
	}

	return token, nil
}

func CheckOTP(ctx context.Context, otp string) (uuid.UUID, error) {
	var userID uuid.UUID

	userIDStr, err := db.RDB.Get(ctx, hashToken(otp)).Result()
	if errors.Is(err, redis.Nil) {
		return uuid.Nil(), errors.New("otp expired or invalid")
	} else if err != nil {
		return userID, fmt.Errorf("get userID: %w", err)
	}

	userID, err = uuid.Parse(userIDStr)
	if err != nil {
		return userID, fmt.Errorf("parse userID: %w", err)
	}

	return userID, nil
}

func DeleteOTP(ctx context.Context, otp string) error {
	err := db.RDB.Del(ctx, hashToken(otp)).Err()
	if err != nil {
		return fmt.Errorf("delete consumed OTP token: %w", err)
	}

	return nil
}
