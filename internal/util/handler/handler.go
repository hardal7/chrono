package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/util/logger"
)

const maxBytes = 256

func Create[T any](f func(http.ResponseWriter, *http.Request, T)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request T

		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Info("Failed to decode JSON")
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			f(w, r, request)
		}
	})
}
