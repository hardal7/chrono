package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/topicevent"
)

func TopicEventRoute(r chi.Router) {
	r.Post("/track", TrackTopicEventHandler)
	r.Get("/all", GetTopicEventsAllHandler)
	r.Get("/today", GetTopicEventsTodayHandler)
}

func TrackTopicEventHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.TrackTopicEventRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = topicevent.Track(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func GetTopicEventsAllHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicEventsRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := topicevent.Get(r.Context(), req)
		processResponse(r.Context(), response{w, resp, err})
	}
}

func GetTopicEventsTodayHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTopicEventsTodayRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := topicevent.GetToday(r.Context(), req)
		processResponse(r.Context(), response{w, resp, err})
	}
}
