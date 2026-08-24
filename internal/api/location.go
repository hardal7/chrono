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
	processRequest(w, r, &req)
	err := location.EditLocationPrivacy(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to edit location privacy", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
