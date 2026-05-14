package topic

import (
	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/util/handler"
)

func Routes(r chi.Router) {
	r.Post("/create", handler.Create(Create, "create topic"))
	r.Post("/edit", handler.Create(Edit, "edit topic"))
	r.Post("/track", handler.Create(Track, "track topic"))
	r.Get("/events", handler.Create(GetEvents, "get topic events"))
}
