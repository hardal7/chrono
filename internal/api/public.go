package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/forms"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

func publicRoutes(r chi.Router) {
	r.Get("/health", PingHandler)
	r.Post("/report", BugReportHandler)
	r.Post("/feature", FeatureRequestHandler)
	r.Post("/register", RegisterUserHandler)
	r.Post("/login", LoginUserHandler)
	r.Post("/password-reset", RequestUserPasswordResetHandler)
	r.Get("/password-reset", UserPasswordResetHandler)
	r.Get(config.AvatarEndpoint+"/{id}", GetUserAvatarHandler)
}

func BugReportHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBugReport
	err := processRequest(w, r, &req)
	if err == nil {
		err = forms.CreateBugReport(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func FeatureRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFeatureRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = forms.CreateFeatureRequest(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = user.Register(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	err := processRequest(w, r, &req)
	if err == nil {
		cookie, err := user.Login(r.Context(), req)
		http.SetCookie(w, &cookie)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func RequestUserPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestUserPasswordResetRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = user.RequestPasswordReset(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func UserPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UserPasswordResetRequest
	err := processRequest(w, r, &req)
	if err != nil {
		otp := r.URL.Query().Get("otp")
		err = user.PasswordReset(r.Context(), otp, req)
		processResponse(r.Context(), response{w, nil, err})
	}
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
		err = dir.Close()
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
		err = file.Close()
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
