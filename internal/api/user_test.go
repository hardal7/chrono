package api_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/hardal7/chrono/internal/api"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/test"
)

func TestMain(m *testing.M) {
	test.Bootstrap()
	api.InitValidator()
	m.Run()
}
func TestRegister(t *testing.T) {
	registerTest.Run(t)
}
func TestLogin(t *testing.T) {
	loginTest.Run(t)
}

var registerTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/register",
	Handler:  api.RegisterUserHandler,
	Cases: []test.Case{
		{
			Name: "successful register",
			Body: dto.RegisterUserRequest{
				Email:    "john2@mail.com",
				Username: "john2doe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name: "duplicate user",
			Body: dto.RegisterUserRequest{
				Email:    "john@mail.com",
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusConflict,
		},
		{
			Name: "missing email",
			Body: dto.RegisterUserRequest{
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "missing username",
			Body: dto.RegisterUserRequest{
				Email:    "john@mail.com",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "missing password",
			Body: dto.RegisterUserRequest{
				Email:    "john@mail.com",
				Username: "johndoe",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:           "empty body",
			Body:           dto.RegisterUserRequest{},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "sql injection attempt",
			Body: dto.RegisterUserRequest{
				Email:    "sqli@mail.com",
				Username: "'; DROP TABLE users;--",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "extremely long input",
			Body: dto.RegisterUserRequest{
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

var loginTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/login",
	Handler:  api.LoginUserHandler,
	Cases: []test.Case{
		{
			Name: "successful login with username",
			Body: dto.LoginUserRequest{
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "successful login with email",
			Body: dto.LoginUserRequest{
				Email:    "john@mail.com",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "incorrect password",
			Body: dto.LoginUserRequest{
				Username: "johndoe",
				Password: "wrongpassword",
			},
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name: "user not found by username",
			Body: dto.LoginUserRequest{
				Username: "unknownuser",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name: "user not found by email",
			Body: dto.LoginUserRequest{
				Email:    "unknown@mail.com",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name: "missing username and email",
			Body: dto.LoginUserRequest{
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name: "empty password",
			Body: dto.LoginUserRequest{
				Username: "johndoe",
				Password: "",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	},
}
