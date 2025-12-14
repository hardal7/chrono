package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hardal7/chrono/internal/config"
	logger "github.com/hardal7/chrono/internal/util"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Authenticating user")

		tokenCookie, err := r.Cookie("Authorization")
		if err == http.ErrNoCookie {
			logger.Info("No token provided")
			logger.Debug(err.Error())
			http.Error(w, "No authorization token found", http.StatusUnauthorized)
			return
		} else if err != nil {
			logger.Info("Failed to get request cookie")
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (any, error) {
				return []byte(config.App.JWT_SECRET), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err == jwt.ErrTokenExpired {
				logger.Info("Token is expired")
				logger.Debug(err.Error())
				http.Error(w, "Token is expired", http.StatusUnauthorized)
				return
			} else if err != nil {
				logger.Info("Invalid Token")
				logger.Debug(err.Error())
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); !ok {
				logger.Info("Failed to parse token")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				userID := int(claims["sub"].(float64))
				logger.Info("Authenticated user")
				ctx := context.WithValue(r.Context(), "userid", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
