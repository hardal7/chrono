package topic

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hardal7/chrono/internal/model"
	logger "github.com/hardal7/chrono/internal/util"

	"github.com/hardal7/chrono/internal/repository"
)

func GetEvents(w http.ResponseWriter, r *http.Request, te model.TopicEventRequest) {
	logger.Info("Getting topic events")
	userID := r.Context().Value("userID").(int)
	topicEvents, err := repository.FindMultiple[model.TopicEvent](r.Context(), "topic_events", "user_id", strconv.Itoa(userID))
	if err != nil {
		logger.Info("Failed to get topic events")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := model.TopicEventResponse{}
		for i := range topicEvents {
			response.Dates[i] = topicEvents[i].Date
			response.TimesTracked[i] = topicEvents[i].TimeTracked
			topic, _ := repository.Get[model.Topic](r.Context(), topicEvents[i].TopicID, "topics")
			response.Topics[i] = topic.Name
		}
		json.NewEncoder(w).Encode(response)
	}
}
