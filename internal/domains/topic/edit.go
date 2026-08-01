package topic

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Edit(w http.ResponseWriter, r *http.Request, er EditRequest) {
	topic, err := Repo.FindUserTopic(r.Context(), r.Context().Value(middleware.UserID).(int), er.Name)
	topic.UpdatedAt = time.Now()
	if err != nil {
		e.ErrNotFound.Handle(w, err, table)
		return
	}
	logger.Info("Editing topic", "topicName", topic.Name)
	if er.DeleteTopic {
		logger.Info("Deleting topic with name: " + topic.Name)
		err := Repo.Delete(r.Context(), topic)
		if err != nil {
			e.ErrDelete.Handle(w, err, table)
			return
		} else {
			logger.Info("Deleted topic")
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	if er.NewName != "" {
		topic.Name = er.NewName
		logger.Info("Changed topic name to " + er.NewName)
	}
	err = Repo.Update(r.Context(), topic)
	if err != nil {
		e.ErrUpdate.Handle(w, err, table)
		return
	} else {
		logger.Info("Changed topic details")
		w.WriteHeader(http.StatusOK)
	}
}
