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

func TestRegisterUser(t *testing.T) {
	registerUserTest.Run(t)
}

func TestLoginUser(t *testing.T) {
	loginUserTest.Run(t)
}

func TestGetUserAccount(t *testing.T) {
	getUserAccountTest.Run(t)
}

const (
	email    = "john@mail.com"
	username = "johndoe"
	password = "strongpassword"
)

var registerUserTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/register",
	Handler:  api.RegisterUserHandler,
	Cases: []test.Case{
		{
			Name: "successful register",
			Body: dto.RegisterUserRequest{
				Email:    "john2@mail.com",
				Username: "john2doe",
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusOK},
		},
		{
			Name: "duplicate user",
			Body: dto.RegisterUserRequest{
				Email:    email,
				Username: username,
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusConflict},
		},
		{
			Name: "missing email",
			Body: dto.RegisterUserRequest{
				Username: username,
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "missing username",
			Body: dto.RegisterUserRequest{
				Email:    email,
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "missing password",
			Body: dto.RegisterUserRequest{
				Email:    email,
				Username: username,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name:             "empty body",
			Body:             dto.RegisterUserRequest{},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "extremely long input",
			Body: dto.RegisterUserRequest{
				Email:    strings.Repeat("a", 5000) + "@mail.com",
				Username: strings.Repeat("u", 5000),
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name:             "malformed json",
			RawBody:          `{"email":email,"username":`,
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
	},
}

var loginUserTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/login",
	Handler:  api.LoginUserHandler,
	Cases: []test.Case{
		{
			Name: "successful login with username",
			Body: dto.LoginUserRequest{
				Username: username,
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusOK},
		},
		{
			Name: "successful login with email",
			Body: dto.LoginUserRequest{
				Email:    email,
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusOK},
		},
		{
			Name: "incorrect password",
			Body: dto.LoginUserRequest{
				Username: username,
				Password: "wrongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "user not found by username",
			Body: dto.LoginUserRequest{
				Username: "unknownuser",
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "user not found by email",
			Body: dto.LoginUserRequest{
				Email:    "unknown@mail.com",
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "missing username and email",
			Body: dto.LoginUserRequest{
				Password: password,
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "empty password",
			Body: dto.LoginUserRequest{
				Username: username,
				Password: "",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
	},
}

var getUserAccountTest = test.Test{
	Method:   http.MethodGet,
	Endpoint: "/user/account",
	Handler:  api.GetUserAccountHandler,
	Cases: []test.Case{
		{
			Name: "successful account details retrieval",
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
				Body:   nil,
			},
		},
	},
}
