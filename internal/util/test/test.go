package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Status int
	Body   any
}

type Case struct {
	Name             string
	Body             any
	RawBody          string
	ExpectedResponse Response
}

type Test struct {
	Method   string
	Endpoint string
	Handler  http.HandlerFunc
	Cases    []Case
}

const testCookie string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODgzMjQzMjEsInN1YiI6MX0.KXelKQAzQNz3mcky05fdI687c_2AkmNN1teGtNA8bzc"

func (test Test) Run(t *testing.T) {
	for _, c := range test.Cases {
		payload, err := c.marshalBody()
		if err != nil {
			logger.Fatal("Failed to marshal test body", err)
		}
		res := httptest.NewRecorder()
		req, err := http.NewRequest(test.Method, test.Endpoint, bytes.NewBuffer(payload))
		req.AddCookie(&http.Cookie{Name: middleware.AuthHeader, Value: testCookie})
		if err != nil {
			logger.Fatal("Failed to create test request", err)
		}
		logger.Info("=== RUNNING TEST ===", "case", c.Name)
		handler := middleware.Authenticate(test.Handler)
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
