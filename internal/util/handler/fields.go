package handler

import (
	"reflect"
)

// TODO: Check for struct tag 'opt' to omit from the return slice
func parseEmptyFields(v any) []string {
	var emptyFields []string
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
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
