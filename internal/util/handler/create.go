package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

const maxBytes = 256

func Create[T any](f func(http.ResponseWriter, *http.Request, T)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request T

		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Warn("Failed to close reader")
			}
		}()
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			ErrBadRequest.Handle(w, nil)
			return
		} else {
			emptyFields := parseEmptyFields(request)
			if len(emptyFields) != 0 {
				msg := "Fields " + strings.Join(emptyFields, ", ") + " cannot be empty"
				err := ErrBadRequest
				err.ExternalInfo = msg
				err.InternalInfo = err.ExternalInfo
				err.Handle(w, nil)
			}
			f(w, r, request)
		}
	})
}

var ErrBadRequest = e.Error{
	InternalInfo: "Bad Request",
	Code:         http.StatusBadRequest,
	ExternalInfo: "Bad Request",
}
