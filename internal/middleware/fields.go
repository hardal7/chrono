package middleware

import (
	"net/http"
	"reflect"
	"strings"

	e "github.com/hardal7/chrono/internal/util/errors"
)

func CheckFields(w http.ResponseWriter, T any) {
	emptyFields := parseEmptyFields(T)
	if len(emptyFields) != 0 {
		msg := "Fields " + strings.Join(emptyFields, ", ") + " cannot be empty"
		err := e.ErrBadRequest
		err.ExternalInfo = msg
		err.InternalInfo = err.ExternalInfo
		err.Handle(w, nil)
	}
}

func parseEmptyFields(v any) []string {
	var emptyFields []string
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		if reflect.TypeOf(v).Field(i).Tag.Get("opt") != "" {
			continue
		}
		field := val.Field(i)
		if field.IsZero() {
			emptyFields = append(emptyFields, reflect.TypeOf(v).Field(i).Name)
		}
	}
	if len(emptyFields) != 0 {
		return emptyFields
	} else {
		return []string{}
	}
}
