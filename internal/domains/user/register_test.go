package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hardal7/chrono/internal/util/handler"
	"github.com/hardal7/chrono/internal/util/test"
)

func TestMain(m *testing.M) {
	test.Bootstrap()
	m.Run()
}

func TestRegister(t *testing.T) {
	product := User{Email: "mail@com", Username: "johndoe", Password: "strongpassword"}
	payload, _ := json.Marshal(product)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal("Failed to make request " + err.Error())
	}

	handler.Create(Register, "register user")(res, req)
	if res.Code != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, res.Code)
	}
}
