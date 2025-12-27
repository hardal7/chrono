package repository

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type CRUDObject struct {
	ID             string
	NumberOfFields int
	FieldNames     []string
	FieldValues    []any
}

func IsDuplicate(ctx context.Context, v any, table string) (bool, error) {
	query := "SELECT COUNT(1) FROM " + table + " WHERE id = $1;"
	var exists int
	err := DB.QueryRow(ctx, query, parseModel(v).ID).Scan(&exists)

	if exists == 0 {
		return false, err
	} else {
		return true, err
	}
}

func Get[T any](ctx context.Context, id int, table string) (T, error) {
	query := "SELECT * FROM " + table + " WHERE id = $1 LIMIT 1;"
	row, err := DB.Query(ctx, query, id)
	model, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	return model, err
}

func Delete(ctx context.Context, v any, table string) error {
	query := "DELETE FROM " + table + " WHERE ID = $1;"
	_, err := DB.Exec(ctx, query, parseModel(v).ID)

	return err
}

func Create(ctx context.Context, v any, table string) error {
	query := "INSERT INTO " + table + createQueryValues(parseModel(v)) + ";"
	_, err := DB.Exec(ctx, query, parseModel(v).FieldValues...)

	return err
}

func Update(ctx context.Context, v any, table string) error {
	query := "UPDATE " + table + " SET " + updateQueryValues(parseModel(v)) + ";"
	_, err := DB.Exec(ctx, query, tidyFields(parseModel(v)).FieldValues...)

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
	// SQL VALUES start from 1 hence i+1 is necessary here to offset
	for i := 0; i < object.NumberOfFields; i++ {
		if i != object.NumberOfFields-1 {
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
	cleanObject.ID = object.ID

	var nonEmptyFields int
	for i := 0; i < object.NumberOfFields; i++ {
		var emptyField bool
		if reflect.ValueOf(object.FieldValues[i]).Kind() == reflect.TypeFor[time.Time]().Kind() {
			if !object.FieldValues[i].(time.Time).IsZero() {
				emptyField = false
			}
		} else if object.FieldValues[i] != "" {
			emptyField = false
		}

		if !emptyField {
			cleanObject.FieldValues[nonEmptyFields] = object.FieldValues[i]
			cleanObject.FieldNames[nonEmptyFields] = object.FieldNames[i]
			nonEmptyFields++
		}
		cleanObject.NumberOfFields = nonEmptyFields
	}

	return cleanObject
}
