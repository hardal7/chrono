package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hardal7/chrono/internal/util/logger"
)

func processRequest(w http.ResponseWriter, r *http.Request, req any) {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}

func processResponse(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", "application/json")

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(resp)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Trace("Returning Response")
	logger.Trace(strings.TrimSpace(buf.String()))

	_, err = w.Write(buf.Bytes())
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
