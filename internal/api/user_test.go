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
func TestEditUserAccount(t *testing.T) {
	editUserAccountTest.Run(t)
}
func TestGetUserAccount(t *testing.T) {
	getUserAccountTest.Run(t)
}

// TODO: Use the middlewares
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
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusCreated},
		},
		{
			Name: "duplicate user",
			Body: dto.RegisterUserRequest{
				Email:    "john@mail.com",
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "missing email",
			Body: dto.RegisterUserRequest{
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "missing username",
			Body: dto.RegisterUserRequest{
				Email:    "john@mail.com",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "missing password",
			Body: dto.RegisterUserRequest{
				Email:    "john@mail.com",
				Username: "johndoe",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name:             "empty body",
			Body:             dto.RegisterUserRequest{},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "sql injection attempt",
			Body: dto.RegisterUserRequest{
				Email:    "sqli@mail.com",
				Username: "'; DROP TABLE users;--",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name: "extremely long input",
			Body: dto.RegisterUserRequest{
				Email:    strings.Repeat("a", 5000) + "@mail.com",
				Username: strings.Repeat("u", 5000),
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
		{
			Name:             "malformed json",
			RawBody:          `{"email":"john@mail.com","username":`,
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
				Username: "johndoe",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusOK},
		},
		{
			Name: "successful login with email",
			Body: dto.LoginUserRequest{
				Email:    "john@mail.com",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusOK},
		},
		{
			Name: "incorrect password",
			Body: dto.LoginUserRequest{
				Username: "johndoe",
				Password: "wrongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusUnauthorized},
		},
		{
			Name: "user not found by username",
			Body: dto.LoginUserRequest{
				Username: "unknownuser",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusNotFound},
		},
		{
			Name: "user not found by email",
			Body: dto.LoginUserRequest{
				Email:    "unknown@mail.com",
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusNotFound},
		},
		{
			Name: "missing username and email",
			Body: dto.LoginUserRequest{
				Password: "strongpassword",
			},
			ExpectedResponse: test.Response{Status: http.StatusNotFound},
		},
		{
			Name: "empty password",
			Body: dto.LoginUserRequest{
				Username: "johndoe",
				Password: "",
			},
			ExpectedResponse: test.Response{Status: http.StatusBadRequest},
		},
	},
}

var editUserAccountTest = test.Test{
	Method:   http.MethodPost,
	Endpoint: "/account",
	Handler:  api.EditUserAccountHandler,
	Cases: []test.Case{
		{
			Name: "successful username change",
			Body: dto.EditUserAccountRequest{
				NewUsername: "johnNewdoe",
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "successful password change",
			Body: dto.EditUserAccountRequest{
				NewPassword: "newPassword123",
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "successful username and password change",
			Body: dto.EditUserAccountRequest{
				NewUsername: "johnNewdoe",
				NewPassword: "newPassword123",
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "successful account deletion",
			Body: dto.EditUserAccountRequest{
				DeleteAccount: true,
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "successful username change and account deletion",
			Body: dto.EditUserAccountRequest{
				NewUsername:   "johnNewdoe",
				DeleteAccount: true,
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "successful password change and account deletion",
			Body: dto.EditUserAccountRequest{
				NewPassword:   "newPassword123",
				DeleteAccount: true,
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "successful username password change and account deletion",
			Body: dto.EditUserAccountRequest{
				NewUsername:   "johnNewdoe",
				NewPassword:   "newPassword123",
				DeleteAccount: true,
			},
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
			},
		},
		{
			Name: "empty request",
			Body: dto.EditUserAccountRequest{},
			ExpectedResponse: test.Response{
				Status: http.StatusBadRequest,
			},
		},
		{
			Name: "empty username",
			Body: dto.EditUserAccountRequest{
				NewUsername: "",
			},
			ExpectedResponse: test.Response{
				Status: http.StatusBadRequest,
			},
		},
		{
			Name: "empty password",
			Body: dto.EditUserAccountRequest{
				NewPassword: "",
			},
			ExpectedResponse: test.Response{
				Status: http.StatusBadRequest,
			},
		},
	},
}

var getUserAccountTest = test.Test{
	Method:   http.MethodGet,
	Endpoint: "/account",
	Handler:  api.EditUserAccountHandler,
	Cases: []test.Case{
		{
			Name: "successful account details retrieval",
			ExpectedResponse: test.Response{
				Status: http.StatusOK,
				Body: dto.GetUserAccountResponse{
					Username: "johndoe",
					Email:    "john@mail.com",
					// CreatedAt: ,
				},
			},
		},
	},
}
