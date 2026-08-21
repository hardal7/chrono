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
const RequestID Key = "requestID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Authenticating user")

		auth := r.Header.Get(AuthHeader)
		if auth == "" {
			logger.Debug("No token provided")
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		auth = strings.TrimPrefix(auth, Bearer)
		token, err := jwt.Parse(auth, func(token *jwt.Token) (any, error) {
			return []byte(config.App.JWT_SECRET), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err == jwt.ErrTokenExpired {
			logger.Debug("Token is expired", err)
			http.Error(w, "Token is expired", http.StatusUnauthorized)
			return
		} else if err != nil {
			logger.Debug("Token has invalid JWT signature", err)
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Debug("Token has invalid claims")
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		sub, ok := claims["sub"].(string)
		if !ok {
			logger.Debug("Token is missing fields")
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		userID, err := uuid.Parse(sub)
		if err != nil {
			logger.Debug("Token does not contain a valid UUID")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		logger.Debug("Authenticated user", "userID", userID.String())
		ctx := context.WithValue(r.Context(), UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
