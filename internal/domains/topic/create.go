package topic

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
)

func Create(w http.ResponseWriter, r *http.Request, cr CreateRequest) {
	logger.Info("Creating topic with name: " + cr.Name)

	topic, err := repository.Find[Topic](r.Context(), "topics", "name", cr.Name)
	if topic.ID != '0' {
		// TODO: Use new error type
		logger.Info("Topic with name " + cr.Name + " already exists")
		http.Error(w, "Topic already exists", http.StatusBadRequest)
		return
	} else if err != nil {
		e.ErrCheckIfDuplicate.Handle(w, err, "topic")
		return
	} else {
		topic := Topic{
			Name:      cr.Name,
			CreatedBy: r.Context().Value(middleware.UserID).(int),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repository.Create(r.Context(), topic, "topics"); err != nil {
			e.ErrCreate.Handle(w, err, "topic")
			return
		} else {
			logger.Info("Created topic: " + cr.Name)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
