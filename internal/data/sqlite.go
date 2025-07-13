package data

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
)

type SQLite struct {
	DB *sql.DB
}

func NewSQLite(ctx context.Context) (*SQLite, error) {
	db, err := sql.Open("sqlite3", filepath.Join(".data", "mocker.db"))
	if err != nil {
		return nil, fmt.Errorf("[func NewSQLite] has failed :: %w", err)
	}

	return &SQLite{
		DB: db,
	}, nil 
}	