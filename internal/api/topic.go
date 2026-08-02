package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/topic"
)

func TopicRoute(r chi.Router) {
	r.Post("/create", CreateTopicHandler)
	r.Post("/edit", EditTopicHandler)
	r.Get("/get", GetTopicHandler)
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTopicRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = topic.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create topic", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func EditTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditTopicRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = topic.Edit(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to edit topic", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	resp, err := topic.Get(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get topic", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
