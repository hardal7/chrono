package topic

import (
	"net/http"

	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
)

func Track(w http.ResponseWriter, r *http.Request, tr TrackTopicRequest) {
	logger.Info("Tracking time for topic with name: " + tr.Topic)
	topic, err := repository.Find[Topic](r.Context(), "topics", "name", tr.Topic)
	if err != nil {
		logger.Info("Topic not found")
		logger.Debug(err.Error())
		http.Error(w, "Topic not found", http.StatusBadRequest)
		return
	} else {
		topicEvent := TopicEvent{
			UserID:      r.Context().Value("userID").(int),
			TopicID:     topic.ID,
			TimeTracked: int(tr.Time.Unix()),
			Date:        int(tr.Date.Unix()),
		}
		if err := repository.Create(r.Context(), topicEvent, "topic_events"); err != nil {
			logger.Info("Failed to track time")
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			logger.Info("Tracked time")
			w.WriteHeader(http.StatusCreated)
		}
	}
}
