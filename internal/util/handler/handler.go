package handler

import (
	"encoding/json"
	"net/http"

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
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			e.ErrDecodeJSON.Handle(w, err)
			return
		} else {
			f(w, r, request)
		}
	})
}
