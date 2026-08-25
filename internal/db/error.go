package db

import "errors"

var ErrRunQuery = errors.New("failed to run query")
var ErrNotFound = errors.New("resource not found")
