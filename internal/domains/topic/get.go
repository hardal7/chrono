package topic

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Get(w http.ResponseWriter, r *http.Request, gr GetRequest) {
	logger.Info("Getting topic", "name", gr.Name)
	userID := r.Context().Value(middleware.UserID).(int)
	t, err := Repo.FindUserTopic(r.Context(), userID, gr.Name)
	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	} else {
		response := GetResponse{
			TotalTime: t.TotalTime,
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
