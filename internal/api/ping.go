package api

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/util/logger"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Got ping")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "service is up",
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
