package user

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/config"
	logger "github.com/hardal7/chrono/internal/util"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hardal7/chrono/internal/model"
	"github.com/hardal7/chrono/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtExpirationDays int = 30
)

func Login(w http.ResponseWriter, r *http.Request, lr model.LoginRequest) {
	logger.Info("Logging user with username: " + lr.Username)

	user, err := repository.GetUserByUsername(r.Context(), lr.Username)
	if err != nil {
		logger.Info("Failed to get user")
		logger.Debug(err.Error())
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		logger.Info("Login attempt with incorect password")
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	} else if err != nil {
		logger.Info("Failed to compare password to hash")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(jwtExpirationDays)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.App.JWT_SECRET))
	if err != nil {
		logger.Info("Failed to generate token")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		cookie := http.Cookie{
			Name:     "Authorization",
			Value:    tokenString,
			Path:     "/",
			MaxAge:   3600 * 24 * jwtExpirationDays,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		logger.Info("Logged user and sent token")
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}
