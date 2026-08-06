package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/logger"
)

func UserRoute(r chi.Router) {
	r.Post("/account", EditUserAccountHandler)
	r.Post("/avatar", UploadUserAvatarHandler)
	r.Get("/avatar/{id}", GetUserAvatar)
	r.Get("/account", GetUserAccountHandler)
	r.With(middleware.Paginate).Get("/top", GetTopUsersHandler)
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

func GetUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := user.GetAccount(r.Context())
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Failed to get account details", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}

func GetTopUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopUsersRequest
	processRequest(w, r, &req)
	resp, err := user.GetTopUsers(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get top users", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}

func EditUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditUserAccountRequest
	processRequest(w, r, &req)
	err := user.EditAccount(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to edit user account", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UploadUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	err := user.UploadAvatar(r.Context(), r.Body)
	if err != nil {
		http.Error(w, "Failed to upload user avatar", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	avatarID := chi.URLParam(r, "id")
	path := filepath.Join("./"+user.AvatarDirectory, avatarID)
	http.ServeFile(w, r, path)
	w.WriteHeader(http.StatusOK)
}
