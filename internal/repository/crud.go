package repository

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	logger "github.com/hardal7/chrono/internal/util"
	"github.com/jackc/pgx/v5"
)

type CRUDObject struct {
	ID             string
	NumberOfFields int
	FieldNames     []string
	FieldValues    []any
}

func IsDuplicate(ctx context.Context, v any, table string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE id = $1;", table)
	logger.Debug("Running query: " + query)
	var exists int
	err := DB.QueryRow(ctx, query, parseModel(v).ID).Scan(&exists)

	if exists == 0 {
		return false, err
	} else {
		return true, err
	}
}

func Find[T any](ctx context.Context, table, field, record string) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE $1 = $2 LIMIT 1;", table)
	row, err := DB.Query(ctx, query, field, record)
	model, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	return model, err
}

func FindMultiple[T any](ctx context.Context, table, field, record string) ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE $1 = $2 LIMIT 1;", table)
	row, err := DB.Query(ctx, query, field, record)
	models, err := pgx.CollectRows(row, pgx.RowToStructByName[T])
	return models, err
}

func Get[T any](ctx context.Context, id int, table string) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 LIMIT 1;", table)
	logger.Debug("Running query: " + query)
	row, err := DB.Query(ctx, query, id)
	model, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	return model, err
}

func Delete(ctx context.Context, v any, table string) error {
	query := fmt.Sprintf("DELETE * FROM %s WHERE id = $1 LIMIT 1;", table)
	logger.Debug("Running query: " + query)
	_, err := DB.Exec(ctx, query, parseModel(v).ID)

	return err
}

func Create(ctx context.Context, v any, table string) error {
	query := "INSERT INTO " + table + createQueryValues(parseModel(v)) + ";"
	logger.Debug("Running query: " + query)
	_, err := DB.Exec(ctx, query, parseModel(v).FieldValues...)

	return err
}

func Update(ctx context.Context, v any, table string) error {
	query := "UPDATE " + table + " SET " + updateQueryValues(parseModel(v)) + ";"
	logger.Debug("Running query: " + query)
	for i := 0; i != tidyFields(parseModel(v)).NumberOfFields; i++ {
		fmt.Printf("%T\n", tidyFields(parseModel(v)).FieldValues[i])
		fmt.Printf("%v\n", tidyFields(parseModel(v)).FieldNames[i])
	}
	request := tidyFields(parseModel(v))
	_, err := DB.Exec(ctx, query, request.FieldValues...)

	return err
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
	var valueString string
	for i := 0; i < object.NumberOfFields; i++ {
		if i == 0 {
			valueString += (" (" + object.FieldNames[i] + ", ")
		} else if i != object.NumberOfFields-1 {
			valueString += (object.FieldNames[i] + ", ")
		} else {
			valueString += (object.FieldNames[i] + ") ")
		}
	}

	valueString += "VALUES"
	for i := 0; i < object.NumberOfFields; i++ {
		if i == 0 {
			// SQL VALUES start from 1 hence i+1 is necessary here to offset
			valueString += (" (" + "$" + strconv.Itoa(i+1) + ", ")
		} else if i != object.NumberOfFields-1 {
			valueString += ("$" + strconv.Itoa(i+1) + ", ")
		} else {
			valueString += ("$" + strconv.Itoa(i+1) + ")")
		}
	}
	return valueString
}

func updateQueryValues(object CRUDObject) string {
	object = tidyFields(object)
	var valueString string
	for i := 0; i < object.NumberOfFields-1; i++ {
		// SQL VALUES start from 1 hence i+1 is necessary here to offset
		if i != object.NumberOfFields-2 {
			valueString += object.FieldNames[i] + " = $" + strconv.Itoa(i+1) + ", "
		} else {
			valueString += object.FieldNames[i] + " = $" + strconv.Itoa(i+1) + " "
		}
	}

	valueString += ("WHERE id = $" + strconv.Itoa(object.NumberOfFields))
	return valueString
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
