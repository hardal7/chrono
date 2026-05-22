package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/repository"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
	"golang.org/x/crypto/bcrypt"
)

func EditAccount(w http.ResponseWriter, r *http.Request, er EditAccountRequest) {
	user, err := repository.Get[User](r.Context(), "users", r.Context().Value(middleware.UserID).(int))
	if err != nil {
		e.ErrNotFound.Handle(w, err, "user")
		return
	} else {
		user.UpdatedAt = time.Now()
		logger.Info("Editing account with username: " + user.Username)
		if er.DeleteAccount {
			logger.Info("Deleting account with username: " + user.Username)
			err := repository.Delete(r.Context(), user, "users")
			if err != nil {
				e.ErrDelete.Handle(w, err, "user")
				return
			} else {
				logger.Info("Deleted account")
				w.WriteHeader(http.StatusOK)
			}
		} else {
			if er.NewUsername != "" {
				user.Username = er.NewUsername
				logger.Info("Changed username to " + er.NewUsername)
			}
			if er.NewPassword != "" {
				logger.Info("Changing account password")
				passwordHash, err := bcrypt.GenerateFromPassword([]byte(er.NewPassword), bcryptCost)
				if err != nil {
					ErrHashPassword.Handle(w, err)
					return
				}
				user.Password = string(passwordHash)
			}
			err := repository.Update(r.Context(), user, "users")
			if err != nil {
				e.ErrUpdate.Handle(w, err, "user")
				return
			} else {
				logger.Info("Changed account details")
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	user, err := repository.Get[User](r.Context(), "users", r.Context().Value(middleware.UserID).(int))
	if err != nil {
		e.ErrNotFound.Handle(w, err, "user")
		return
	}
	response, err := json.Marshal(user)
	if err != nil {
		e.ErrMarshalJSON.Handle(w, err)
		return
	} else {
		logger.Info("Sent account details")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
