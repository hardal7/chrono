package errors

import "net/http"

var ErrNotFound = ErrorWithResource{
	Error: Error{
		InternalInfo: "Failed to find resource: ",
		Code:         http.StatusNotFound,
		ExternalInfo: "User not found",
	},
}
var ErrCheckIfDuplicate = ErrorWithResource{
	Error: Error{
		InternalInfo: "Failed to check if resource is duplicate: ",
		Code:         http.StatusInternalServerError,
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
