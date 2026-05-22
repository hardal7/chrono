package topicevent

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
)

func GetEvents(w http.ResponseWriter, r *http.Request, te TopicEventRequest) {
	logger.Info("Getting topic events")
	userID := r.Context().Value(middleware.UserID).(int)
	topicEvents, err := repository.FindMultiple[TopicEvent](r.Context(), "topic_events", "user_id", strconv.Itoa(userID))
	if err != nil {
		e.ErrNotFound.Handle(w, err, "topic_events")
		return
	} else {
		response := TopicEventResponse{}
		for i := range topicEvents {
			response.Dates[i] = topicEvents[i].Date
			response.TimesTracked[i] = topicEvents[i].TimeTracked
			topic, _ := repository.Get[topic.Topic](r.Context(), "topics", topicEvents[i].TopicID)
			response.Topics[i] = topic.Name
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
