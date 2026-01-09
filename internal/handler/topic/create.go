package topic

import (
	"net/http"
	"time"

	logger "github.com/hardal7/chrono/internal/util"

	"github.com/hardal7/chrono/internal/model"
	"github.com/hardal7/chrono/internal/repository"
)

func Create(w http.ResponseWriter, r *http.Request, tr model.CreateTopicRequest) {
	user, _ := repository.Get[model.User](r.Context(), r.Context().Value("userID").(int), "users")
	logger.Info("Creating topic with name: " + tr.Name)

	topic := model.Topic{
		Name:      tr.Name,
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isDuplicate, err := repository.IsDuplicate(r.Context(), topic, "topics")
	if err != nil {
		logger.Info("Failed to check if topic" + topic.Name + " is duplicate")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isDuplicate {
		logger.Info("Topic with name " + topic.Name + " already exists")
		http.Error(w, "Topic already exists", http.StatusBadRequest)
		return
	} else {
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
