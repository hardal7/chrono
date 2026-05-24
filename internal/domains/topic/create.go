package topic

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Create(w http.ResponseWriter, r *http.Request, cr CreateRequest) {
	logger.Info("Creating topic with name: " + cr.Name)

	topic, err := Repo.FindByName(r.Context(), cr.Name)
	if topic.ID != 0 {
		e.ErrAlreadyExists.Handle(w, err, table)
		return
	} else if err != nil {
		e.ErrCheckIfDuplicate.Handle(w, err, table)
		return
	} else {
		topic := Topic{
			Name:      cr.Name,
			CreatedBy: r.Context().Value(middleware.UserID).(int),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := Repo.Create(r.Context(), topic); err != nil {
			e.ErrCreate.Handle(w, err, table)
			return
		} else {
			logger.Info("Created topic: " + cr.Name)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
