package db

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

// TODO: Reformat this code
func Get[T any](ctx context.Context, table, field, record string) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 LIMIT 1;", table, field)
	logger.Debug("Running query: " + query)
	row, _ := DB.Query(ctx, query, record)
	model, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	return model, err
}

func GetMultiple[T any](ctx context.Context, table, field, record string) ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 LIMIT 1;", table, field)
	logger.Debug("Running query: " + query)
	row, _ := DB.Query(ctx, query, record)
	models, err := pgx.CollectRows(row, pgx.RowToStructByName[T])
	return models, err
}

func Delete(ctx context.Context, table string, v any) error {
	query := fmt.Sprintf("DELETE * FROM %s WHERE id = $1 LIMIT 1;", table)
	logger.Debug("Running query: " + query)
	_, err := DB.Exec(ctx, query, parseModel(v).ID)

	return err
}

func Create(ctx context.Context, table string, v any) error {
	query := fmt.Sprintf("INSERT INTO %s %s;", table, buildCreateString(parseModel(v)))
	logger.Debug("Running query: " + query)
	_, err := DB.Exec(ctx, query, parseModel(v).FieldValues...)

	return err
}

func Update(ctx context.Context, table string, v any) error {
	query := fmt.Sprintf("UPDATE %s SET %s;", table, buildUpdateString(parseModel(v)))
	logger.Debug("Running query: " + query)
	for i := 0; i != tidyFields(parseModel(v)).NumberOfFields; i++ {
		fmt.Printf("%T\n", tidyFields(parseModel(v)).FieldValues[i])
		fmt.Printf("%v\n", tidyFields(parseModel(v)).FieldNames[i])
	}
	request := tidyFields(parseModel(v))
	_, err := DB.Exec(ctx, query, request.FieldValues...)

	return err
}

type sqlObject struct {
	ID             string
	NumberOfFields int
	FieldNames     []string
	FieldValues    []any
}

// Takes a generic model and returns an sqlObject
func parseModel(v any) sqlObject {
	var object sqlObject
	structType := reflect.TypeOf(v)
	structValue := reflect.ValueOf(v)
	object.ID = strconv.FormatInt(structValue.Field(0).Int(), 10)
	object.NumberOfFields = structType.NumField() - 1
	object.FieldNames = make([]string, object.NumberOfFields)
	object.FieldValues = make([]any, object.NumberOfFields)

	for i := 1; i < structType.NumField(); i++ {
		object.FieldValues[i-1] = structValue.Field(i).Interface()
		object.FieldNames[i-1] = string([]byte(structValue.Type().Field(i).Tag.Get("db")))
	}
	return object
}

func buildCreateString(object sqlObject) string {
	var valueString strings.Builder
	for i := range object.NumberOfFields {
		if i == 0 {
			valueString.WriteString(("(" + object.FieldNames[i] + ", "))
		} else if i != object.NumberOfFields-1 {
			valueString.WriteString((object.FieldNames[i] + ", "))
		} else {
			valueString.WriteString((object.FieldNames[i] + ") "))
		}
	}

	valueString.WriteString("VALUES")
	for i := range object.NumberOfFields {
		if i == 0 {
			// SQL VALUES start from 1 hence i+1 is necessary here to offset
			valueString.WriteString((" (" + "$" + strconv.Itoa(i+1) + ", "))
		} else if i != object.NumberOfFields-1 {
			valueString.WriteString(("$" + strconv.Itoa(i+1) + ", "))
		} else {
			valueString.WriteString(("$" + strconv.Itoa(i+1) + ")"))
		}
	}
	return valueString.String()
}

func buildUpdateString(object sqlObject) string {
	object = tidyFields(object)
	var valueString strings.Builder
	for i := range object.NumberOfFields - 1 {
		// SQL VALUES start from 1 hence i+1 is necessary here to offset
		if i != object.NumberOfFields-2 {
			valueString.WriteString(object.FieldNames[i] + " = $" + strconv.Itoa(i+1) + ", ")
		} else {
			valueString.WriteString(object.FieldNames[i] + " = $" + strconv.Itoa(i+1) + " ")
		}
	}

	valueString.WriteString(("WHERE id = $" + strconv.Itoa(object.NumberOfFields)))
	return valueString.String()
}

// Cleans fields of an sqlObject where where values are uninitialized
func tidyFields(object sqlObject) sqlObject {
	var cleanObject sqlObject

	nonEmptyFields := 0
	for i := range object.NumberOfFields {
		emptyField := false
		// For fields with value types time.Time, the uninitialized value is different than 0
		if reflect.ValueOf(object.FieldValues[i]).Kind() == reflect.TypeFor[time.Time]().Kind() {
			if object.FieldValues[i].(time.Time).IsZero() {
				emptyField = true
			}
		} else if object.FieldValues[i] == "" {
			emptyField = true
		}

		if !emptyField {
			cleanObject.FieldValues = append(cleanObject.FieldValues, object.FieldValues[i])
			cleanObject.FieldNames = append(cleanObject.FieldNames, object.FieldNames[i])
			nonEmptyFields++
		}
	}

	cleanObject.ID = object.ID
	cleanObject.FieldValues = append(cleanObject.FieldValues, object.ID)
	cleanObject.FieldNames = append(cleanObject.FieldNames, "id")
	nonEmptyFields++

	cleanObject.NumberOfFields = nonEmptyFields

	return cleanObject
}
