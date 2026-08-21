package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/service/user"
)

func UserRoute(r chi.Router) {
	r.Post("/account", EditUserAccountHandler)
	r.Post("/avatar", UploadUserAvatarHandler)
	r.Get("/account", GetUserAccountHandler)
	r.Get("/profile/{username}", GetUserProfileHandler)
	r.With(middleware.Paginate).Get("/top", GetTopUsersHandler)
}

func GetUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := user.GetAccount(r.Context())
	if err != nil {
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

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	resp, err := user.GetProfile(r.Context(), username)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}
