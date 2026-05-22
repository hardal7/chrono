package user

import (
	"net/http"

	e "github.com/hardal7/chrono/internal/util/errors"
)

var ErrIncorrectPassword = e.Error{
	InternalInfo: "Login attempt with incorect password",
	Code:         http.StatusUnauthorized,
	ExternalInfo: "Incorrect Password",
}
var ErrCompareHash = e.Error{
	InternalInfo: "Failed to compare password to hash",
	Code:         http.StatusInternalServerError,
}
var ErrGenerateToken = e.Error{
	InternalInfo: "Failed to generate token",
	Code:         http.StatusInternalServerError,
}
var ErrHashPassword = e.Error{
	InternalInfo: "Failed to hash password",
	Code:         http.StatusInternalServerError,
}
var ErrAlreadyRegistered = e.Error{
	InternalInfo: "User is already registered",
	Code:         http.StatusInternalServerError,
	ExternalInfo: "User exists",
}
