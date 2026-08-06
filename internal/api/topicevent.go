package api

import (
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
	processRequest(w, r, &req)
	err := topicevent.Track(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to track topic event", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetTopicEventsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicEventsRequest
	processRequest(w, r, &req)
	resp, err := topicevent.Get(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get topic events", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}

func GetTopicEventsTodayHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := topicevent.GetToday(r.Context())
	if err != nil {
		http.Error(w, "Failed to get topic events today", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}
