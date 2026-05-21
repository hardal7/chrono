package user

import (
	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/util/handler"
)

func Routes(r chi.Router) {
	r.Post("/account", handler.Create(EditAccount))
	r.Get("/account", GetAccount)
}
