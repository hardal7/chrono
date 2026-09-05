package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/service/user"
)

func UserRoute(r chi.Router) {
	r.Post("/activity", PostUserActivityHandler)
	r.Get("/account", GetUserAccountHandler)
	r.Post("/account", EditUserAccountHandler)
	r.Delete("/account", DeleteUserAccountHandler)
	r.Post("/logout", LogoutUserHandler)
	r.Post("/avatar", UploadUserAvatarHandler)
	r.Delete("/avatar", DeleteUserAvatarHandler)
	r.Get("/profile/{username}", GetUserProfileHandler)
	r.With(middleware.Paginate).Get("/top", GetTopUsersHandler)
}

func PostUserActivityHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func EditUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditUserAccountRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = user.EditAccount(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func GetUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := user.GetAccount(r.Context())
	processResponse(r.Context(), response{w, resp, err})
}

func DeleteUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	err := user.DeleteAccount(r.Context())
	processResponse(r.Context(), response{w, nil, err})
}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	err := user.Logout(r.Context())
	processResponse(r.Context(), response{w, nil, err})
}

func UploadUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	err := user.UploadAvatar(r.Context(), r.Body)
	processResponse(r.Context(), response{w, nil, err})
}

func DeleteUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	err := user.DeleteAvatar(r.Context())
	processResponse(r.Context(), response{w, nil, err})
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	resp, err := user.GetProfile(r.Context(), username)
	processResponse(r.Context(), response{w, resp, err})
}

func GetTopUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopUsersRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := user.GetTopUsers(r.Context(), req)
		processResponse(r.Context(), response{w, resp, err})
	}
}
