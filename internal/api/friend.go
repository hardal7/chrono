package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/friend"
)

func FriendRoute(r chi.Router) {
	r.Post("/create", CreateFriendRequestHandler)
	r.Post("/accept", AcceptFriendRequestHandler)
	r.Delete("/remove", RemoveFriendHandler)
	r.Get("/all", GetAllFriendRequestsHandler)
}

func CreateFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFriendRequestRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = friend.CreateRequest(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}

func AcceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AcceptFriendRequestRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = friend.AcceptRequest(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}

func RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RemoveFriendRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = friend.Remove(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}

func GetAllFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := friend.GetAll(r.Context())
	processResponse(response{w, resp}, err)
}
