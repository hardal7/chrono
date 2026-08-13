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
}

func CreateFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFriendRequestRequest
	processRequest(w, r, &req)
	err := friend.CreateRequest(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create friend request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func AcceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AcceptFriendRequestRequest
	processRequest(w, r, &req)
	err := friend.AcceptRequest(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to accept friend request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RemoveFriendRequest
	processRequest(w, r, &req)
	err := friend.Remove(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to remove friend", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
