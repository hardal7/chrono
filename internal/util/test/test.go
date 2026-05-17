package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/stretchr/testify/assert"
)

type Case struct {
	Name           string
	Body           any
	RawBody        string
	ExpectedStatus int
}

type Test struct {
	Method   string
	Endpoint string
	Handler  http.HandlerFunc
	Cases    []Case
}

func (test Test) Run(t *testing.T) {
	for _, c := range test.Cases {
		payload, err := c.marshalBody()
		if err != nil {
			logger.Error("Failed to marshal test body")
		}
		res := httptest.NewRecorder()
		req, err := http.NewRequest(test.Method, test.Endpoint, bytes.NewBuffer(payload))
		if err != nil {
			logger.Error("Failed to create test request")
		}
		test.Handler(res, req)
		assert.Equal(t, c.ExpectedStatus, res.Code, "Test case %q failed", c.Name)
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
