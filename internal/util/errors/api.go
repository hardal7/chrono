package errors

import "net/http"

var ErrBadRequest = Error{
	InternalInfo: "Bad Request",
	Code:         http.StatusBadRequest,
	ExternalInfo: "Bad Request",
}
