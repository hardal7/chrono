package topicevent

import (
	"net/http"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Track(w http.ResponseWriter, r *http.Request, tr TrackTopicRequest) {
	logger.Info("Tracking time", "topicName", tr.Topic)
	topic, err := topic.Repo.FindUserTopic(r.Context(), r.Context().Value(middleware.UserID).(int), tr.Topic)
	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	} else {
		topicEvent := TopicEvent{
			UserID:      r.Context().Value(middleware.UserID).(int),
			TopicID:     topic.ID,
			TimeTracked: tr.TimeSeconds,
			Date:        tr.Date,
		}
		if err := Repo.Create(r.Context(), topicEvent); err != nil {
			e.ErrCreate.Handle(w, err, table)
			return
		} else {
			logger.Info("Tracked time for topic: " + tr.Topic)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
