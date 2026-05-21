package topic

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/repository"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Edit(w http.ResponseWriter, r *http.Request, er EditRequest) {
	topic, err := repository.Find[Topic](r.Context(), "topics", "name", er.Name)
	topic.UpdatedAt = time.Now()
	if err != nil {
		logger.Info("Failed to get topic")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		logger.Info("Editing topic with name: " + topic.Name)
		if er.DeleteTopic {
			logger.Info("Deleting topic with name: " + topic.Name)
			err := repository.Delete(r.Context(), topic, "topics")
			if err != nil {
				logger.Info("Failed to delete topic")
				logger.Debug(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				logger.Info("Deleted topic")
				w.WriteHeader(http.StatusOK)
			}
		} else {
			if er.NewName != "" {
				topic.Name = er.NewName
				logger.Info("Changed topic name to " + er.NewName)
			}
			err := repository.Update(r.Context(), topic, "topics")
			if err != nil {
				logger.Info("Failed to change topic details")
				logger.Debug(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				logger.Info("Changed topic details")
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
