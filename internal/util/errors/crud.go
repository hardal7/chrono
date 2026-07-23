package errors

import "net/http"

var ErrNotFound = ErrorWithResource{
	Error: Error{
		InternalInfo: "Failed to find resource in table: ",
		Code:         http.StatusNotFound,
		ExternalInfo: "Not found",
	},
}
var ErrAlreadyExists = ErrorWithResource{
	Error: Error{
		InternalInfo: "Resource already exists: ",
		Code:         http.StatusConflict,
		ExternalInfo: "Already exists",
	},
}
var ErrCreate = ErrorWithResource{
	Error: Error{
		InternalInfo: "Failed to create resource: ",
		Code:         http.StatusInternalServerError,
	},
}
var ErrDelete = ErrorWithResource{
	Error: Error{
		InternalInfo: "Failed to delete resource: ",
		Code:         http.StatusInternalServerError,
	},
}
var ErrUpdate = ErrorWithResource{
	Error: Error{
		InternalInfo: "Failed to update resource: ",
		Code:         http.StatusInternalServerError,
	},
}
