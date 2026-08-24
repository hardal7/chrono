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
	err := processRequest(w, r, &req)
	if err == nil {
		err = topic.Create(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}

func EditTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditTopicRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = topic.Edit(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}

func GetAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := topic.GetAll(r.Context())
	processResponse(response{w, resp}, err)
}

func GetNamedTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicNamedRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := topic.GetNamed(r.Context(), req)
		processResponse(response{w, resp}, err)
	}
}
