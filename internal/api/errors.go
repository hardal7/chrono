package api

import (
	"errors"
	"net/http"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/apierror"
	"github.com/hardal7/chrono/internal/util/logger"
)

func handleErrors(err error, w http.ResponseWriter) {
	if errors.Is(err, db.ErrRunQuery) {
		logger.Warn(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, db.ErrBeginTransaction) {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, db.ErrCommitTransaction) {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, db.ErrNotFound) {
		logger.Debug(err.Error())
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if errors.Is(err, apierror.ErrAlreadyExists) {
		logger.Debug(err.Error())
		http.Error(w, "Already Exists", http.StatusConflict)
		return
	}

	if errors.Is(err, apierror.ErrUnauthorized) {
		logger.Debug(err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		logger.Debug(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
