package topic

import (
	"net/http"

	"github.com/hardal7/chrono/internal/model"
	"github.com/hardal7/chrono/internal/repository"
	logger "github.com/hardal7/chrono/internal/util"
)

func EditTopic(w http.ResponseWriter, r *http.Request, tr model.EditTopicRequest) {
	topic, err := repository.GetTopicByName(r.Context(), tr.Name)
	if err != nil {
		logger.Info("Failed to get topic")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		logger.Info("Editing topic with name: " + topic.Name)
		if tr.DeleteTopic {
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
			if tr.NewName != "" {
				topic.Name = tr.NewName
				logger.Info("Changed topic name to " + tr.NewName)
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
