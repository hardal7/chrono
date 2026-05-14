package session

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/repository"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Create(w http.ResponseWriter, r *http.Request, csr CreateSessionRequest) {
	logger.Info("Creating session with name: " + csr.Name)

	session := Session{
		Name:      csr.Name,
		Password:  csr.Password,
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
