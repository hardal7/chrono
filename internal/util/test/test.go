package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Method   string
	Endpoint string
	Handler  http.HandlerFunc
	Cases    []Case
}

type Case struct {
	Name             string
	Body             any
	RawBody          string
	ExpectedResponse Response
}

type Response struct {
	Status int
	Body   any
}

func (test Test) Run(t *testing.T) {
	for _, c := range test.Cases {
		payload, err := c.marshalBody()
		if err != nil {
			logger.Fatal("Failed to marshal test body", err)
		}
		res := httptest.NewRecorder()
		req, err := http.NewRequest(test.Method, test.Endpoint, bytes.NewBuffer(payload))
		if err != nil {
			logger.Fatal("Failed to create test request", err)
		}
		testUUID, _ := uuid.Parse("b60aa148-0849-4246-8fbd-3e7500316989")
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserID, testUUID))
		handler := middleware.LogRequest(test.Handler)
		logger.Info("=== RUNNING TEST ===", "case", c.Name)
		handler.ServeHTTP(res, req)
		assert.Equal(t, c.ExpectedResponse.Status, res.Code, "Test case %q failed", c.Name)
	}
}

func (c Case) marshalBody() ([]byte, error) {
	if c.RawBody != "" {
		return []byte(c.RawBody), nil
	}

	if c.Body != nil {
		return json.Marshal(c.Body)
	}

	return []byte{}, nil
}
