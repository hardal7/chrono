package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/topicevent"
)

func TopicEventRoute(r chi.Router) {
	r.Post("/track", TrackTopicEventHandler)
	r.Get("/get", GetTopicEventsHandler)
	r.Get("/today", GetTopicEventsTodayHandler)
}

func TrackTopicEventHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.TrackTopicEventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = topicevent.Track(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to track topic event", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetTopicEventsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicEventsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	resp, err := topicevent.Get(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get topic events", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetTopicEventsTodayHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := topicevent.GetToday(r.Context())
	if err != nil {
		http.Error(w, "Failed to get topic events today", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
