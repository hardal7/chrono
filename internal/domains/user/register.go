package user

import (
	"net/http"
	"time"

	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost int = 10

func Register(w http.ResponseWriter, r *http.Request, rr RegisterRequest) {
	logger.Info("Registering user", "username", rr.Username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcryptCost)
	if err != nil {
		ErrHashPassword.Handle(w, err)
		return
	}
	_, err = Repo.FindByUsername(r.Context(), rr.Username)
	if err != pgx.ErrNoRows {
		e.ErrAlreadyExists.Handle(w, err, table)
		return
	} else {
		user := User{
			Email:     rr.Email,
			Username:  rr.Username,
			Password:  string(passwordHash),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := Repo.Create(r.Context(), user); err != nil {
			e.ErrCreate.Handle(w, err, table)
			return
		} else {
			logger.Info("Registered user: " + rr.Username)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
