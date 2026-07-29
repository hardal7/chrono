package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

const maxBytes = 256

func Create[T any](f func(http.ResponseWriter, *http.Request, T)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			e.ErrBadRequest.Handle(w, nil)
			logger.Debug(err.Error())
			return
		}
		logger.Trace(string(body))
		r.Body = io.NopCloser(bytes.NewReader(body))

		var request T
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Warn("Failed to close reader")
			}
		}()
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			e.ErrBadRequest.Handle(w, nil)
			return
		} else {
			middleware.CheckFields(w, request)
			f(w, r, request)
		}
	})
}
