package errors

import "net/http"

var ErrMarshalJSON = Error{
	InternalInfo: "Failed to marshal JSON",
	Code:         http.StatusInternalServerError,
}
var ErrDecodeJSON = Error{
	InternalInfo: "Failed to decode JSON",
	Code:         http.StatusInternalServerError,
}
var ErrEncodeJSON = Error{
	InternalInfo: "Failed to encode JSON",
	Code:         http.StatusInternalServerError,
}
