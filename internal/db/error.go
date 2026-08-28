package db

import "errors"

var ErrRunQuery = errors.New("failed to run query")
var ErrBeginTransaction = errors.New("failed to begin transaction")
var ErrNotFound = errors.New("resource not found")
