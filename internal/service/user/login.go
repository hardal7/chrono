package user

import (
	"context"
	"net/http"
	"time"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const jwtExpirationDays int = 30

func Login(ctx context.Context, r dto.LoginUserRequest) (http.Cookie, error) {
	logger.Debug("Logging user", "username", r.Username)

	var err error
	var u db.User
	if r.Username != "" {
		u, err = conn.Queries.GetUserByUsername(ctx, r.Username)
	} else if r.Email != "" {
		u, err = conn.Queries.GetUserByEmail(ctx, r.Email)
	}
	if err != nil {
		logger.Debug("Failed to get user", err)
		return http.Cookie{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		logger.Debug("Wrong password", err)
		return http.Cookie{}, err
	} else if err != nil {
		logger.Debug("Failed to hash password", err)
		return http.Cookie{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(jwtExpirationDays)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.App.JWT_SECRET))
	if err != nil {
		logger.Debug("Failed to sign token", err)
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:     middleware.AuthHeader,
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600 * 24 * jwtExpirationDays,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	logger.Debug("Logged user and sent token", "username", r.Username)
	return cookie, nil
}
