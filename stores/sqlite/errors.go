// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package sqlite

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrCreateSchema        = Error("create schema")
	ErrDatabaseExists      = Error("database exists")
	ErrForeignKeysDisabled = Error("foreign keys disabled")
	ErrInvalidPath         = Error("invalid path")
	ErrMissingAdminSecret  = Error("missing admin secret")
	ErrNotDirectory        = Error("not a directory")
	ErrPragmaReturnedNil   = Error("pragma returned nil")
)
