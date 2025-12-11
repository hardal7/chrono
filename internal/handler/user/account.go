package user

import (
	"net/http"
	"time"

	"github.com/hardal7/study/internal/model"
	"github.com/hardal7/study/internal/repository"
	logger "github.com/hardal7/study/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func EditAccount(w http.ResponseWriter, r *http.Request, er model.EditAccountRequest) {
	user, err := repository.GetUserByID(r.Context(), r.Context().Value("userid").(int))
	if err != nil {
		logger.Info("Failed to get user")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		if er.DeleteAccount {
			logger.Info("Deleting account with username" + user.Username)
			err := repository.DeleteUser(r.Context(), user)
			if err != nil {
				logger.Info("Failed to delete account")
				logger.Debug(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				logger.Info("Deleted account")
				w.WriteHeader(http.StatusOK)
			}
		} else {
			user.UpdatedAt = time.Now()
			if er.NewUsername != "" {
				user.Username = er.NewUsername
				logger.Info("Changing account username from " + user.Username + " to " + er.NewUsername)
			}
			if er.NewPassword != "" {
				logger.Info("Changing account password")
				passwordHash, err := bcrypt.GenerateFromPassword([]byte(er.NewPassword), bcryptCost)
				if err != nil {
					logger.Info("Could not hash password")
					logger.Debug(err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				user.Password = string(passwordHash)
			}
			err := repository.UpdateUser(r.Context(), user)
			if err != nil {
				logger.Info("Failed to change account details")
				logger.Debug(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				logger.Info("Changed account details")
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
