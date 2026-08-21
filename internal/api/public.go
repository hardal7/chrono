package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/config"
)

func publicRoutes(r chi.Router) {
	r.Post("/register", RegisterUserHandler)
	r.Post("/login", LoginUserHandler)
	r.Get("/health", PingHandler)
	r.Get(config.AvatarEndpoint+"/{id}", GetUserAvatarHandler)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	processRequest(w, r, &req)
	err := user.Register(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	processRequest(w, r, &req)
	cookie, err := user.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to login user", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func GetUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	avatarID := chi.URLParam(r, "id")
	path := filepath.Join(user.AvatarDirectory, avatarID)
	http.ServeFile(w, r, path)
}
