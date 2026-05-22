package topicevent

import (
	"net/http"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
)

func Track(w http.ResponseWriter, r *http.Request, tr TrackTopicRequest) {
	logger.Info("Tracking time for topic with name: " + tr.Topic)
	topic, err := repository.Find[topic.Topic](r.Context(), "topics", "name", tr.Topic)
	if err != nil {
		e.ErrNotFound.Handle(w, err, "topic_events")
		return
	} else {
		topicEvent := TopicEvent{
			UserID:      r.Context().Value(middleware.UserID).(int),
			TopicID:     topic.ID,
			TimeTracked: int(tr.Time.Unix()),
			Date:        int(tr.Date.Unix()),
		}
		if err := repository.Create(r.Context(), topicEvent, "topic_events"); err != nil {
			e.ErrCreate.Handle(w, err, "topic_events")
			return
		} else {
			logger.Info("Tracked time")
			w.WriteHeader(http.StatusCreated)
		}
	}
}
