package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hardal7/chrono/internal/db"
	e "github.com/hardal7/chrono/internal/util/error"
	"github.com/hardal7/chrono/internal/util/logger"
)

func processRequest(w http.ResponseWriter, r *http.Request, req any) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}

	err = validate.Struct(req)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	return nil
}

type response struct {
	w    http.ResponseWriter
	body any
}

func processResponse(r response, err error) {
	if errors.Is(err, db.ErrRunQuery) {
		logger.Error(err.Error())
		http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, e.ErrAlreadyExists) {
		logger.Error(err.Error())
		http.Error(r.w, "Already Exists", http.StatusConflict)
		return
	}

	if err != nil {
		logger.Error(err.Error())
		http.Error(r.w, "Bad Request", http.StatusBadRequest)
		return
	}

	if r.body != nil {
		r.w.Header().Set("Content-Type", "application/json")

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(r.body)
		if err != nil {
			logger.Error(err.Error())
			http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Trace("Returning Response")
		logger.Trace(strings.TrimSpace(buf.String()))

		_, err = r.w.Write(buf.Bytes())
		if err != nil {
			logger.Error(err.Error())
			http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		r.w.WriteHeader(http.StatusOK)
	}

}
