package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hardal7/chrono/internal/util/config"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

type Key string

const AuthHeader string = "Authorization"
const UserID Key = "userID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Authenticating user")

		tokenCookie, err := r.Cookie(AuthHeader)
		if err == http.ErrNoCookie {
			ErrNoToken.Handle(w, err)
			return
		} else if err != nil {
			ErrRetrieveToken.Handle(w, err)
			return
		} else {
			token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (any, error) {
				return []byte(config.App.JWT_SECRET), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err == jwt.ErrTokenExpired {
				ErrExpiredToken.Handle(w, err)
				return
			} else if err != nil {
				ErrInvalidToken.Handle(w, err)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); !ok {
				ErrParseToken.Handle(w, err)
				return
			} else {
				userID := int(claims["sub"].(float64))
				logger.Info("Authenticated user")
				ctx := context.WithValue(r.Context(), UserID, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}

var ErrNoToken = e.Error{
	InternalInfo: "No token provided",
	Code:         http.StatusUnauthorized,
	ExternalInfo: "No authorization token found",
}
var ErrRetrieveToken = e.Error{
	InternalInfo: "Failed to retrieve authorization token",
	Code:         http.StatusInternalServerError,
}
var ErrExpiredToken = e.Error{
	InternalInfo: "Expired token provided",
	Code:         http.StatusUnauthorized,
	ExternalInfo: "Token is expired",
}
var ErrInvalidToken = e.Error{
	InternalInfo: "Invalid token provided",
	Code:         http.StatusUnauthorized,
	ExternalInfo: "Token is invalid",
}
var ErrParseToken = e.Error{
	InternalInfo: "Failed to parse token",
	Code:         http.StatusInternalServerError,
}
