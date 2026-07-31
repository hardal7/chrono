package topicevent

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Get(w http.ResponseWriter, r *http.Request, gr GetRequest) {
	// TODO: Don't return all events but only specified from the request "gr"
	logger.Info("Getting topic events")
	userID := r.Context().Value(middleware.UserID).(int)
	topicEvents, err := Repo.GetAll(r.Context(), userID)
	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	} else {
		response := GetResponse{}
		for i := range topicEvents {
			response.Dates[i] = topicEvents[i].Date
			response.TimesTracked[i] = topicEvents[i].TimeTracked
			t, _ := topic.Repo.FindByID(r.Context(), topicEvents[i].TopicID)
			response.Topics[i] = t.Name
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
