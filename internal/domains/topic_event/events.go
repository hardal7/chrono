package topicevent

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetEvents(w http.ResponseWriter, r *http.Request, te TopicEventRequest) {
	logger.Info("Getting topic events")
	userID := r.Context().Value(middleware.UserID).(int)
	topicEvents, err := Repo.GetAll(r.Context(), userID)
	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	} else {
		response := TopicEventResponse{}
		for i := range topicEvents {
			response.Dates[i] = topicEvents[i].Date
			response.TimesTracked[i] = topicEvents[i].TimeTracked
			topic, _ := topic.Repo.FindByID(r.Context(), topicEvents[i].TopicID)
			response.Topics[i] = topic.Name
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
