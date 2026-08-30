package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
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
	err  error
}

func processResponse(ctx context.Context, r response) {
	if r.err != nil {
		handleErrors(r.err, r.w)
		return
	}

	if r.body != nil {
		r.w.Header().Set("Content-Type", "application/json")

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(r.body)
		if err != nil {
			logger.Debug(err.Error())
			http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Debug("Returning Response")
		requestID := ctx.Value(requestctx.RequestID).(string)
		logger.Debug(strings.TrimSpace(buf.String()), "requestID", requestID)

		_, err = r.w.Write(buf.Bytes())
		if err != nil {
			logger.Debug(err.Error())
			http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		r.w.WriteHeader(http.StatusOK)
	}
}
