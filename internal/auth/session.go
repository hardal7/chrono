package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"uuid"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/redis/go-redis/v9"
)

const sessionExpiration = time.Hour * 24 * 30

func GenerateCookie(ctx context.Context) (http.Cookie, error) {
	var cookie http.Cookie

	token, err := generateSession(ctx)
	if err != nil {
		return cookie, err
	}

	cookie = http.Cookie{
		Name:     AuthCookie,
		Value:    token,
		Path:     "/api",
		MaxAge:   int(sessionExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie, nil
}

func generateSession(ctx context.Context) (string, error) {
	token, err := generateToken()
	if err != nil {
		return token, fmt.Errorf("generate session: %w", err)
	}

	err = db.RDB.Set(ctx, hashToken(token), requestctx.GetUserID(ctx).String(), sessionExpiration).Err()
	if err != nil {
		return token, fmt.Errorf("set session values: %w", err)
	}

	return token, nil
}

func checkSession(ctx context.Context, session string) (uuid.UUID, error) {
	var userID uuid.UUID

	userIDStr, err := db.RDB.Get(ctx, hashToken(session)).Result()
	if errors.Is(err, redis.Nil) {
		return uuid.Nil(), errors.New("session expired or invalid")
	} else if err != nil {
		return userID, fmt.Errorf("get userID: %w", err)
	}

	userID, err = uuid.Parse(userIDStr)
	if err != nil {
		return userID, fmt.Errorf("parse userID: %w", err)
	}

	return userID, nil
}

func DeleteSession(ctx context.Context) error {
	err := db.RDB.Del(ctx, hashToken(requestctx.GetSessionID(ctx))).Err()
	if err != nil {
		return fmt.Errorf("delete consumed Session token: %w", err)
	}

	return nil
}
