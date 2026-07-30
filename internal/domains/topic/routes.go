package topic

import (
	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/util/handler"
)

func Routes(r chi.Router) {
	r.Get("/get", handler.Create(Get))
	r.Post("/create", handler.Create(Create))
	r.Post("/edit", handler.Create(Edit))
}
