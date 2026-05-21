package topicevent

import (
	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/util/handler"
)

func Routes(r chi.Router) {
	r.Post("/track", handler.Create(Track))
	r.Get("/events", handler.Create(GetEvents))
}
