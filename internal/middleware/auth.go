package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

const AuthHeader string = "Authorization"
const Bearer string = "Bearer "

type Key string

const UserID Key = "userID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Authenticating user")

		auth := r.Header.Get(AuthHeader)
		if auth == "" {
			logger.Error("No token provided")
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		} else {
			auth = strings.TrimPrefix(auth, Bearer)
			token, err := jwt.Parse(auth, func(token *jwt.Token) (any, error) {
				return []byte(config.App.JWT_SECRET), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err == jwt.ErrTokenExpired {
				logger.Error("Token is expired")
				http.Error(w, "Token is expired", http.StatusUnauthorized)
				return
			} else if err != nil {
				logger.Debug(err.Error())
				logger.Error("Token is invalid")
				http.Error(w, "Token is invalid", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); !ok {
				logger.Error("Token is invalid")
				http.Error(w, "Token is invalid", http.StatusInternalServerError)
				return
			} else {
				userID, err := uuid.Parse(claims["sub"].(string))
				if err != nil {
					logger.Debug(err.Error())
					logger.Error("Failed to parse token")
					http.Error(w, "Failed to parse token", http.StatusInternalServerError)
					return
				}
				logger.Info("Authenticated user")
				ctx := context.WithValue(r.Context(), UserID, userID)

				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
