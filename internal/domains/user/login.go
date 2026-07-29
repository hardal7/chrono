package user

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const jwtExpirationDays int = 30

func Login(w http.ResponseWriter, r *http.Request, lr LoginRequest) {
	logger.Info("Logging user", "username", lr.Username)

	var err error
	var user User
	if lr.Username != "" {
		user, err = Repo.FindByUsername(r.Context(), lr.Username)
	} else if lr.Email != "" {
		user, err = Repo.FindByEmail(r.Context(), lr.Email)
	}

	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		ErrIncorrectPassword.Handle(w, err)
		return
	} else if err != nil {
		ErrCompareHash.Handle(w, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(jwtExpirationDays)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.App.JWT_SECRET))
	if err != nil {
		ErrGenerateToken.Handle(w, err)
		return
	} else {
		cookie := http.Cookie{
			Name:     middleware.AuthHeader,
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
