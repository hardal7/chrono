package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hardal7/chrono/internal/util/config"
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
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Failed to retrieve token", http.StatusInternalServerError)
			return
		} else {
			token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (any, error) {
				return []byte(config.App.JWT_SECRET), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err == jwt.ErrTokenExpired {
				http.Error(w, "Token is expired", http.StatusUnauthorized)
				return
			} else if err != nil {
				http.Error(w, "Token is invalid", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); !ok {
				http.Error(w, "Failed to parse token", http.StatusInternalServerError)
				return
			} else {
				userID := int32(claims["sub"].(float64))
				logger.Info("Authenticated user")
				ctx := context.WithValue(r.Context(), UserID, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
