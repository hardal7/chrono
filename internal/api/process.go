package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/apierror"
	"github.com/hardal7/chrono/internal/util/logger"
)

func processRequest(w http.ResponseWriter, r *http.Request, req any) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}

	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	return nil
}

type response struct {
	w    http.ResponseWriter
	body any
	err  error
}

func processResponse(ctx context.Context, r response) {
	if errors.Is(r.err, db.ErrRunQuery) {
		logger.Debug(r.err.Error())
		http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(r.err, db.ErrNotFound) {
		logger.Debug(r.err.Error())
		http.Error(r.w, "Not Found", http.StatusNotFound)
		return
	}

	if errors.Is(r.err, apierror.ErrAlreadyExists) {
		logger.Debug(r.err.Error())
		http.Error(r.w, "Already Exists", http.StatusConflict)
		return
	}

	if r.err != nil {
		logger.Debug(r.err.Error())
		http.Error(r.w, "Bad Request", http.StatusBadRequest)
		return
	}

	if r.body != nil {
		r.w.Header().Set("Content-Type", "application/json")

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(r.body)
		if err != nil {
			logger.Debug(err.Error())
			http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Debug("Returning Response")
		requestID := ctx.Value(middleware.RequestID).(string)
		logger.Debug(strings.TrimSpace(buf.String()), "requestID", requestID)

		_, err = r.w.Write(buf.Bytes())
		if err != nil {
			logger.Debug(err.Error())
			http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		r.w.WriteHeader(http.StatusOK)
	}
}
