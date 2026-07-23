package user

import (
	"net/http"
	"strings"
	"testing"

	"github.com/hardal7/chrono/internal/util/handler"
	"github.com/hardal7/chrono/internal/util/test"
)

func TestMain(m *testing.M) {
	test.Bootstrap()
	m.Run()
}

func TestRegister(t *testing.T) {
	registerTest.Run(t)
}

var registerTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/register",
	Handler:  handler.Create(Register),
	Cases: []test.Case{
		{
			Name: "successful register",
			Body: User{
				Email:    "john2@mail.com",
				Username: "john2doe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name: "duplicate user",
			Body: User{
				Email:    "john@mail.com",
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusConflict,
		},
		{
			Name: "missing email",
			Body: User{
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "missing username",
			Body: User{
				Email:    "john@mail.com",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "missing password",
			Body: User{
				Email:    "john@mail.com",
				Username: "johndoe",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:           "empty body",
			Body:           User{},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "sql injection attempt",
			Body: User{
				Email:    "sqli@mail.com",
				Username: "'; DROP TABLE users;--",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "extremely long input",
			Body: User{
				Email:    strings.Repeat("a", 5000) + "@mail.com",
				Username: strings.Repeat("u", 5000),
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:           "malformed json",
			RawBody:        `{"email":"john@mail.com","username":`,
			ExpectedStatus: http.StatusBadRequest,
		},
	},
}
