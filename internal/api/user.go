package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/user"
	"github.com/hardal7/chrono/internal/util/logger"
)

func UserRoute(r chi.Router) {
	r.Post("/account", EditUserAccountHandler)
	r.Post("/avatar", UploadUserAvatarHandler)
	r.Get("/account", GetUserAccountHandler)
	r.Get("/top", GetTopUsersHandler)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = user.Register(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetTopUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	resp, err := user.GetTopUsers(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get top users", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func EditUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditUserAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = user.EditAccount(r.Context(), req)
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
}
