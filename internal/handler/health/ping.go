package health

import (
	"encoding/json"
	"net/http"

	logger "github.com/hardal7/chrono/internal/util"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	logger.Info("Got ping")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
	})
}
