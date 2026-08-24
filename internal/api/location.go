package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/service/location"
)

func LocationRoute(r chi.Router) {
	r.Post("/privacy", EditLocationPrivacyHandler)
}

func EditLocationPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.EditLocationPrivacyRequest
	err := processRequest(w, r, &req)
	if err == nil {
		err = location.EditLocationPrivacy(r.Context(), req)
		processResponse(response{w, nil}, err)
	}
}
