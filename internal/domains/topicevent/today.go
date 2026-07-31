package topicevent

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetToday(w http.ResponseWriter, r *http.Request) {
	logger.Info("Getting total time tracked today")
	userID := r.Context().Value(middleware.UserID).(int)
	topicEvents, err := Repo.GetToday(r.Context(), userID)
	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	} else {
		var totalTime int
		for i := range topicEvents {
			totalTime += topicEvents[i].TimeTracked
		}
		w.Header().Set("Content-Type", "application/json")
		response := GetTodayResponse{
			TotalTime: totalTime,
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
