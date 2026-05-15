package topic

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
)

func Create(w http.ResponseWriter, r *http.Request, tr CreateTopicRequest) {
	logger.Info("Creating topic with name: " + tr.Name)

	topic, err := repository.Find[Topic](r.Context(), "topics", "name", tr.Name)
	if topic.ID != '0' {
		logger.Info("Topic with name " + tr.Name + " already exists")
		http.Error(w, "Topic already exists", http.StatusBadRequest)
		return
	} else if err != nil {
		logger.Info("Failed to check if topic" + topic.Name + " is duplicate")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		topic := Topic{
			Name:      tr.Name,
			CreatedBy: r.Context().Value("userID").(int),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repository.Create(r.Context(), topic, "topics"); err != nil {
			logger.Info("Failed to create topic with name: " + tr.Name)
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			logger.Info("Created topic: " + tr.Name)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
