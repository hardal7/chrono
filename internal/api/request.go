package api

import (
	"encoding/json"
	"net/http"

	logger "github.com/hardal7/chrono/internal/util"
)

func CreateRequest[T any](f func(http.ResponseWriter, *http.Request, T), operation ...string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request T
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Info("Failed to " + operation[0])
			logger.Info("Failed to decode JSON")
			logger.Debug(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			f(w, r, request)
		}
	})
}
