package topic

import (
	"net/http"
	"time"

	logger "github.com/hardal7/chrono/internal/util"

	"github.com/hardal7/chrono/internal/model"
	"github.com/hardal7/chrono/internal/repository"
)

func Create(w http.ResponseWriter, r *http.Request, tr model.CreateTopicRequest) {
	logger.Info("Registering topic with name: " + tr.Name)

	topic := model.Topic{
		Name:      tr.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isDuplicate, err := repository.IsDuplicate(r.Context(), topic, "topics")
	if err != nil {
		logger.Info("Failed to check if topic " + topic.Name + " is duplicate")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isDuplicate {
		logger.Info("Topic " + tr.Name + " exists")
		http.Error(w, "Topic exists", http.StatusBadRequest)
		return
	} else {
		if err := repository.Create(r.Context(), topic, "topics"); err != nil {
			logger.Info("Failed to create topic: " + tr.Name)
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			logger.Info("Created Topic: " + tr.Name)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
