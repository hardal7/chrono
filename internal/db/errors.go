package db

import "errors"

var (
	ErrRunQuery          = errors.New("failed to run query")
	ErrNotFound          = errors.New("resource not found")
	ErrBeginTransaction  = errors.New("failed to begin transaction")
	ErrCommitTransaction = errors.New("failed to commit transaction")
)
