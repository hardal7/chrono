package errors

import (
	"net/http"

	"github.com/hardal7/chrono/internal/util/logger"
)

type Error struct {
	InternalInfo string
	Code         int
	ExternalInfo string
}

type ErrorWithResource struct {
	Error    Error
	Resource string
}

func (e Error) Handle(w http.ResponseWriter, debug error) {
	if e.Code == http.StatusInternalServerError {
		e.ExternalInfo = "Internal Server Error"
		logger.Warn(e.InternalInfo)
	} else {
		logger.Error(e.InternalInfo)
	}
	if debug != nil {
		logger.Debug(debug.Error())
	}
	http.Error(w, e.ExternalInfo, e.Code)
}

func (e ErrorWithResource) Handle(w http.ResponseWriter, debug error, resource string) {
	e.Resource = resource
	e.Error.InternalInfo = e.Error.InternalInfo + e.Resource
	e.Error.Handle(w, debug)
}
