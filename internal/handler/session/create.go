package session

import (
	"net/http"
	"time"

	logger "github.com/hardal7/chrono/internal/util"

	"github.com/hardal7/chrono/internal/model"
	"github.com/hardal7/chrono/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost int = 10
)

func Create(w http.ResponseWriter, r *http.Request, csr model.CreateSessionRequest) {
	logger.Info("Creating session with name: " + csr.Name)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(csr.Password), bcryptCost)
	if err != nil {
		logger.Info("Failed to create session: Could not hash password")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session := model.Session{
		Name:     csr.Name,
		Password: string(passwordHash),
		// TODO: Use auth middleware beforehand to get userid
		// Admin: ,
		Expiry:    time.Now().Add(time.Minute * time.Duration(csr.Expiry)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isDuplicate, err := repository.IsDuplicate(r.Context(), session, "sessions")
	if err != nil {
		logger.Info("Failed to check if session " + csr.Name + " is duplicate")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isDuplicate {
		logger.Info("Session with name " + session.Name + " is already created")
		http.Error(w, "Session with name "+session.Name+" is already created", http.StatusBadRequest)
		return
	} else {
		if err := repository.Create(r.Context(), session, "sessions"); err != nil {
			logger.Info("Failed to create session: " + csr.Name)
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			logger.Info("Created session: " + csr.Name)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
