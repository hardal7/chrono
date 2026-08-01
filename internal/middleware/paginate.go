package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	e "github.com/hardal7/chrono/internal/util/errors"
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
			e.ErrBadRequest.Handle(w, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var req listRequest
		if err := json.Unmarshal(body, &req); err != nil {
			e.ErrBadRequest.Handle(w, err)
			return
		}

		if req.Limit > 20 || req.Limit < 0 {
			e.ErrBadRequest.Handle(w, "Invalid limit range")
			return
		}

		next.ServeHTTP(w, r)
	})
}
