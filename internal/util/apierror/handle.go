package apierror

import (
	"context"
	"errors"
	"net/http"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func Handle(ctx context.Context, w http.ResponseWriter, err error) {
	e := logger.Err(err).With("requestID", requestctx.GetRequestID(ctx))
	msg := "Request Failed"

	if errors.Is(err, db.ErrRunQuery) {
		e.Warn(msg)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, db.ErrBeginTransaction) {
		e.Error(msg)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, db.ErrCommitTransaction) {
		e.Error(msg)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, db.ErrNotFound) {
		e.Debug(msg)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if errors.Is(err, ErrAlreadyExists) {
		e.Debug(msg)
		http.Error(w, "Already Exists", http.StatusConflict)
		return
	}

	if errors.Is(err, ErrUnauthorized) {
		e.Debug(msg)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		e.Debug(msg)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
