package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/service/user"
)

func UserRoute(r chi.Router) {
	r.Get("/account", GetUserAccountHandler)
	r.Post("/account", EditUserAccountHandler)
	r.Post("/avatar", UploadUserAvatarHandler)
	r.Get("/profile/{username}", GetUserProfileHandler)
	r.With(middleware.Paginate).Get("/top", GetTopUsersHandler)
}

func GetUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := user.GetAccount(r.Context())
	processResponse(response{w, resp}, err)
}

func GetTopUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopUsersRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := user.GetTopUsers(r.Context(), req)
		processResponse(response{w, resp}, err)
	}
}

func EditUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditUserAccountRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = user.EditAccount(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}

func UploadUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	err := user.UploadAvatar(r.Context(), r.Body)
	processResponse(response{w, nil}, err)
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	resp, err := user.GetProfile(r.Context(), username)
	processResponse(response{w, resp}, err)
}
