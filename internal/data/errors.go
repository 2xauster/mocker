package data

import (
	"database/sql"
	"errors"

	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

func SQLiteErrorComparator(err error) error {
	if err == nil {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
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

func RedisErrorComparator(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return errs.NewError(err, errs.RedisErrorType, errs.ErrNotFound)
	}

	var redisErr redis.Error
	if errors.As(err, &redisErr) {
		return errs.NewError(err, errs.RedisErrorType, errs.ErrUndefined)
	}

	return err
}