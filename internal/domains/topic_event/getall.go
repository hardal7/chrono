package topicevent

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
)

func GetAll(w http.ResponseWriter, r *http.Request, te TopicEventRequest) {
	logger.Info("Getting topic events")
	userID := r.Context().Value(middleware.UserID).(int)
	topicEvents, err := repository.FindMultiple[TopicEvent](r.Context(), "topic_events", "user_id", strconv.Itoa(userID))
	if err != nil {
		logger.Info("Failed to get topic events")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := TopicEventResponse{}
		for i := range topicEvents {
			response.Dates[i] = topicEvents[i].Date
			response.TimesTracked[i] = topicEvents[i].TimeTracked
			topic, _ := repository.Get[topic.Topic](r.Context(), topicEvents[i].TopicID, "topics")
			response.Topics[i] = topic.Name
		}
		json.NewEncoder(w).Encode(response)
	}
}
