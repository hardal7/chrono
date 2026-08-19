package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/util/config"
)

func publicRoutes(r chi.Router) {
	r.Post("/register", RegisterUserHandler)
	r.Post("/login", LoginUserHandler)
	r.Get("/health", PingHandler)
	r.Get(config.AvatarEndpoint+"/{id}", GetUserAvatarHandler)
}
