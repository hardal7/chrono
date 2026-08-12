package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/topic"
	"github.com/hardal7/chrono/internal/util/logger"
)

func TopicRoute(r chi.Router) {
	r.Post("/create", CreateTopicHandler)
	r.Post("/edit", EditTopicHandler)
	r.Get("/get", GetTopicHandler)
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTopicRequest
	processRequest(w, r, &req)
	err := topic.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create topic", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func EditTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditTopicRequest
	processRequest(w, r, &req)
	err := topic.Edit(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to edit topic", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicRequest
	processRequest(w, r, &req)
	resp, err := topic.Get(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get topics", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
