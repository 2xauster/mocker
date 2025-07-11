package data_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	_ "github.com/mattn/go-sqlite3"
)

func TestPrepareStatement(t *testing.T) {
	fields := map[string]data.SQLField{
		"ID": {
			Name:        "ID",
			Datatype:    "TEXT",
			Constraints: "PRIMARY KEY UNIQUE NOT NULL",
		},
		"Name": {
			Name:        "Name",
			Datatype:    "TEXT",
			Constraints: "NOT NULL",
		},
	}

	stmt := data.PrepareCreateTableStmt("user", fields)
	t.Logf("Generated SQL: %s", stmt)

	sqlLower := strings.ToLower(stmt)

	must := []string{
		"create", "table", "if", "not", "exists",
		"user", "(", "id", "text", "primary", "key", "unique", "not", "null",
		"name", "text", "not", "null", ")", ";",
	}

	for _, token := range must {
		if !strings.Contains(sqlLower, token) {
			t.Errorf("expected SQL to contain %q, but it did not", token)
		}
	}
}

func TestCreateTable(t *testing.T) {
	db, err := sql.Open("sqlite3", filepath.Join("..", "..", ".data", "test.db"))
	if err != nil {
		t.Fatal("failed to open db:", err)
	}
	defer db.Close()

	user := TestUserEntity{}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := data.CreateTable(ctx, db, user)
	t.Logf("Executing :: %s", stmt)

	if err != nil {
		t.Fatal("failed to create table:", err)
	}

	rows, err := db.Query("PRAGMA table_info(TestUserEntity);")
	if err != nil {
		t.Fatal("failed to inspect table:", err)
	}
	defer rows.Close()

	var foundCols []string
	for rows.Next() {
		var (
			cid        int
			name       string
			typ        string
			notnull    int
			dfltValue  sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			t.Fatal("failed to scan row:", err)
		}
		t.Logf("Column: %s, Type: %s, NotNull: %t, PK: %t", name, typ, notnull != 0, pk != 0)
		foundCols = append(foundCols, name)
	}

	expected := []string{"ID", "Name"}
	for _, col := range expected {
		found := false
		for _, actual := range foundCols {
			if actual == col {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected column %q not found in table", col)
		}
	}
}