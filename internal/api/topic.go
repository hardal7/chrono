package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/topic"
)

func TopicRoute(r chi.Router) {
	r.Post("/create", CreateTopicHandler)
	r.Post("/edit", EditTopicHandler)
	r.Get("/all", GetAllTopicsHandler)
	r.Get("/named", GetNamedTopicHandler)
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
	w.WriteHeader(http.StatusOK)
}

func GetAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := topic.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get all topics", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}

func GetNamedTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicNamedRequest
	processRequest(w, r, &req)
	resp, err := topic.GetNamed(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get named topic", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}
