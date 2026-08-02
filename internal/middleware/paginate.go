package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hardal7/chrono/internal/util/logger"
)

type listRequest struct {
	Cursor int `json:"cursor"`
	Limit  int `json:"limit"`
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Paginating Request")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var req listRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.Limit > 20 || req.Limit < 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
