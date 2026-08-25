package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/jackc/pgx/v5"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	day               = 3600 * 24
	jwtExpirationDays = 30
)

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(jwtExpirationDays)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.App.JWT_SECRET))
	if err != nil {
		return cookie, fmt.Errorf("Failed to sign token: %w", err)
	}

	cookie = http.Cookie{
		Name:     middleware.AuthHeader,
		Value:    tokenString,
		Path:     "/api",
		MaxAge:   jwtExpirationDays,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie, nil
}
