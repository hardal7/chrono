package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
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

	// SECURITY: https://go.dev/blog/osroot
	// Default avatars are stored as symlinks: e.g. /avatars/00uu-00uu -> /avatars/default/3
	// Resolved paths need to be handled with care to secure from path traversal attacks
	dir, err := os.OpenRoot(user.AvatarDirectory)
	if err != nil {
		logger.Warn("Failed to open avatar directory", err)
		http.NotFound(w, r)
		return
	}
	defer func() {
		err := dir.Close()
		if err != nil {
			logger.Warn("Failed to close avatar directory")
		}
	}()

	file, err := dir.Open(avatarID)
	if err != nil {
		logger.Debug("Failed to open avatar file", err)
		http.NotFound(w, r)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Warn("Failed to close avatar file")
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.IsDir() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}
