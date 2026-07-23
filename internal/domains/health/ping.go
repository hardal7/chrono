package health

import (
	"encoding/json"
	"net/http"

	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	logger.Info("Got ping")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "service is up",
	}); err != nil {
		e.ErrMarshalJSON.Handle(w, err)
	}
}
