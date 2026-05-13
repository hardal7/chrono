package repository

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

func IsDuplicate(ctx context.Context, v any, table string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE id = %s;", table, parseModel(v).ID)
	logger.Debug("Running query: " + query)
	var exists int
	err := DB.QueryRow(ctx, query).Scan(&exists)

	if exists == 0 {
		return false, err
	} else {
		return true, err
	}
}

func Find[T any](ctx context.Context, table, field, record string) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s LIMIT 1;", table, field, record)
	row, err := DB.Query(ctx, query)
	model, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	return model, err
}

func FindMultiple[T any](ctx context.Context, table, field, record string) ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s LIMIT 1;", table, field, record)
	row, err := DB.Query(ctx, query)
	models, err := pgx.CollectRows(row, pgx.RowToStructByName[T])
	return models, err
}

func Get[T any](ctx context.Context, id int, table string) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %s LIMIT 1;", table, strconv.Itoa(id))
	logger.Debug("Running query: " + query)
	row, err := DB.Query(ctx, query)
	model, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	return model, err
}

func Delete(ctx context.Context, v any, table string) error {
	query := fmt.Sprintf("DELETE * FROM %s WHERE id = %s LIMIT 1;", table, parseModel(v).ID)
	logger.Debug("Running query: " + query)
	_, err := DB.Exec(ctx, query)

	return err
}

func Create(ctx context.Context, v any, table string) error {
	query := fmt.Sprintf("INSERT INTO %s %s;", table, createQueryValues(parseModel(v)))
	logger.Debug("Running query: " + query)
	_, err := DB.Exec(ctx, query, parseModel(v).FieldValues...)

	return err
}

func Update(ctx context.Context, v any, table string) error {
	query := fmt.Sprintf("UPDATE %s SET %s;", table, updateQueryValues(parseModel(v)))
	logger.Debug("Running query: " + query)
	for i := 0; i != tidyFields(parseModel(v)).NumberOfFields; i++ {
		fmt.Printf("%T\n", tidyFields(parseModel(v)).FieldValues[i])
		fmt.Printf("%v\n", tidyFields(parseModel(v)).FieldNames[i])
	}
	request := tidyFields(parseModel(v))
	_, err := DB.Exec(ctx, query, request.FieldValues...)

	return err
}

type CRUDObject struct {
	ID             string
	NumberOfFields int
	FieldNames     []string
	FieldValues    []any
}

func parseModel(v any) CRUDObject {
	var object CRUDObject
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

func createQueryValues(object CRUDObject) string {
	var valueString strings.Builder
	for i := 0; i < object.NumberOfFields; i++ {
		if i == 0 {
			valueString.WriteString(("(" + object.FieldNames[i] + ", "))
		} else if i != object.NumberOfFields-1 {
			valueString.WriteString((object.FieldNames[i] + ", "))
		} else {
			valueString.WriteString((object.FieldNames[i] + ") "))
		}
	}

	valueString.WriteString("VALUES")
	for i := 0; i < object.NumberOfFields; i++ {
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

func updateQueryValues(object CRUDObject) string {
	object = tidyFields(object)
	var valueString strings.Builder
	for i := 0; i < object.NumberOfFields-1; i++ {
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

func tidyFields(object CRUDObject) CRUDObject {
	var cleanObject CRUDObject

	nonEmptyFields := 0
	for i := 0; i < object.NumberOfFields; i++ {
		emptyField := false
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
