package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

const tokenExpiration = time.Hour * 24 * 30

func Login(ctx context.Context, r dto.LoginUserRequest) (http.Cookie, error) {
	var err error
	var u query.User
	var cookie http.Cookie

	if r.Username != "" {
		u, err = db.Queries.GetUserByUsername(ctx, r.Username)
	} else if r.Email != "" {
		u, err = db.Queries.GetUserByEmail(ctx, r.Email)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return cookie, fmt.Errorf("User not found")
	} else if err != nil {
		return cookie, fmt.Errorf("Failed to get user: %w: %w", db.ErrRunQuery, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return cookie, fmt.Errorf("Wrong password: %w", err)
	} else if err != nil {
		return cookie, fmt.Errorf("Failed to hash password: %w", err)
	}

	token, err := auth.GenerateToken()
	if err != nil {
		return cookie, fmt.Errorf("Failed to generate token: %q", err)
	}
	hashedToken := auth.HashToken(token, []byte(config.App.HashSecret))

	err = db.Queries.CreateSessionToken(ctx, query.CreateSessionTokenParams{
		UserID: u.ID,
		Expiry: time.Now().Add(tokenExpiration),
		Hash:   hashedToken,
	})
	if err != nil {
		return cookie, fmt.Errorf("Failed to create session token: %w: %w", db.ErrRunQuery, err)
	}

	cookie = http.Cookie{
		Name:     auth.AuthCookie,
		Value:    token,
		Path:     "/api",
		MaxAge:   int(tokenExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie, nil
}
