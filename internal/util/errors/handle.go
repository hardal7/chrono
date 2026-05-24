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

func (e Error) Error() string {
	return e.ExternalInfo
}

type ErrorWithResource struct {
	Error    Error
	Resource string
}

func (e Error) Handle(w http.ResponseWriter, debug any) {
	if e.Code == http.StatusInternalServerError {
		e.ExternalInfo = "Internal Server Error"
		logger.Warn(e.InternalInfo)
	} else {
		logger.Error(e.InternalInfo)
	}

	if debug != nil {
		switch v := debug.(type) {
		case error:
			logger.Debug(v.Error())
		case string:
			logger.Debug(v)
		}
	}
	http.Error(w, e.ExternalInfo, e.Code)
}

func (e ErrorWithResource) Handle(w http.ResponseWriter, debug error, resource string) {
	e.Resource = resource
	e.Error.InternalInfo = e.Error.InternalInfo + e.Resource
	e.Error.Handle(w, debug)
}
