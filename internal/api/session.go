package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/session"
)

func SessionRoute(r chi.Router) {
	r.Post("/create", CreateSessionHandler)
	r.Post("/edit", EditSessionHandler)
	r.Post("/join", JoinSessionHandler)
	r.Get("/named", GetNamedSessionHandler)
	r.Get("/all", GetAllSessionsHandler)
}

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSessionRequest
	processRequest(w, r, &req)
	err := session.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func EditSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditSessionRequest
	processRequest(w, r, &req)
	err := session.Edit(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to edit session", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func JoinSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.JoinSessionRequest
	processRequest(w, r, &req)
	err := session.Join(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to join session", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetNamedSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetSessionNamedRequest
	resp, err := session.GetNamed(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get named session", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}
func GetAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := session.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get all sessions", http.StatusBadRequest)
		return
	}
	processResponse(w, resp)
}
