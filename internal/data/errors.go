package data

import (
	"database/sql"
	"errors"

	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/mattn/go-sqlite3"
)

func SQLiteErrorComparator(err error) error {
	if err == sql.ErrNoRows {
		return errs.NewError(err, errs.SQLErrorType, errs.ErrNotFound)
	}
	
	var sqliteErr sqlite3.Error
	if !errors.As(err, &sqliteErr) {
		return err
	}

	switch sqliteErr.ExtendedCode {
	case sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey:
		return errs.NewError(err, errs.SQLErrorType, errs.ErrAlreadyExists)
	default:
		return err
	}
}