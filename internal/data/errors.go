package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

const (
	ErrNotFound = iota
	ErrAlreadyExists
	ErrInvalidSyntax
)

// Errors for SQL. So sweet it will give you diabetes.
type SQLError struct {
	Err  error
	Code int
}

func NewSQLError(err error, code int) SQLError {
	return SQLError{
		Err:  err,
		Code: code,
	}
}

func (e SQLError) Error() string {
	return fmt.Sprintf("sql error (%d): %v", e.Code, e.Err)
}


func (e SQLError) Unwrap() error {
	return e.Err
}

func (e SQLError) Is(target error) bool {
	t, ok := target.(SQLError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func SugarifyErrors(err error) error {
	if err == sql.ErrNoRows {
		return NewSQLError(err, ErrNotFound)
	}
	
	var sqliteErr sqlite3.Error 
	if !errors.As(err, &sqliteErr) {
		return err
	}

	switch sqliteErr.ExtendedCode {
	case sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey:
		return NewSQLError(err, ErrAlreadyExists)
	default:
		return err
	}
}