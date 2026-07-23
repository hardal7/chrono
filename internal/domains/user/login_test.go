package user

import (
	"net/http"
	"testing"

	"github.com/hardal7/chrono/internal/util/handler"
	"github.com/hardal7/chrono/internal/util/test"
)

func TestLogin(t *testing.T) {
	loginTest.Run(t)
}

var loginTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/login",
	Handler:  handler.Create(Login),
	Cases: []test.Case{
		{
			Name: "successful login with username",
			Body: LoginRequest{
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "successful login with email",
			Body: LoginRequest{
				Email:    "john@mail.com",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name: "incorrect password",
			Body: LoginRequest{
				Username: "johndoe",
				Password: "wrongpassword",
			},
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name: "user not found by username",
			Body: LoginRequest{
				Username: "unknownuser",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name: "user not found by email",
			Body: LoginRequest{
				Email:    "unknown@mail.com",
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name: "missing username and email",
			Body: LoginRequest{
				Password: "strongpassword",
			},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name: "empty password",
			Body: LoginRequest{
				Username: "johndoe",
				Password: "",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	},
}
