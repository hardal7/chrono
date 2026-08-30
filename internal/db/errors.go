package db

import "errors"

var ErrRunQuery = errors.New("failed to run query")
var ErrNotFound = errors.New("resource not found")
var ErrBeginTransaction = errors.New("failed to begin transaction")
var ErrCommitTransaction = errors.New("failed to commit transaction")
