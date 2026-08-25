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
	err := processRequest(w, r, &req)
	if err == nil {
		err = session.Create(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func EditSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditSessionRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = session.Edit(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func JoinSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.JoinSessionRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = session.Join(r.Context(), req)
		processResponse(r.Context(), response{w, nil, err})
	}
}

func GetNamedSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetSessionNamedRequest
	err := processRequest(w, r, &req)
	if err == nil {
		resp, err := session.GetNamed(r.Context(), req)
		processResponse(r.Context(), response{w, resp, err})
	}
}
func GetAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := session.GetAll(r.Context())
	processResponse(r.Context(), response{w, resp, err})
}
