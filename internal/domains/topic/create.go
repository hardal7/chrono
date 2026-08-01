package topic

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

func Create(w http.ResponseWriter, r *http.Request, cr CreateRequest) {
	logger.Info("Creating topic", "topicName", cr.Name)

	_, err := Repo.FindUserTopic(r.Context(), r.Context().Value(middleware.UserID).(int), cr.Name)
	if err != pgx.ErrNoRows {
		e.ErrAlreadyExists.Handle(w, err, table)
		return
	}
	topic := Topic{
		Name:      cr.Name,
		TotalTime: 0,
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
