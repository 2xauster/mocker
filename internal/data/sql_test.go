package data_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func createTestDB() (*sql.DB, error) {
	testDBPath := filepath.Join("..", "..", ".data", "unit_test.db")

	if err := os.MkdirAll(filepath.Dir(testDBPath), os.ModePerm); err != nil {
		return nil, err
	}

	if _, err := os.Stat(testDBPath); os.IsNotExist(err) {
		file, err := os.Create(testDBPath)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	return sql.Open("sqlite3", testDBPath)
}

func dropTestTable(db *sql.DB, tableName string) {
	stmt := fmt.Sprintf(`DROP TABLE "%s"`, tableName)
	_, err := db.Exec(stmt)
	if err != nil {
	} // eat five star, do nothing
}

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
	db, err := createTestDB()
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
	defer dropTestTable(db, "TestUserEntity")

	rows, err := db.Query("PRAGMA table_info(TestUserEntity);")
	if err != nil {
		t.Fatal("failed to inspect table:", err)
	}
	defer rows.Close()

	var foundCols []string
	for rows.Next() {
		var (
			cid       int
			name      string
			typ       string
			notnull   int
			dfltValue sql.NullString
			pk        int
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

func TestInsert(t *testing.T) {
	db, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS insertUsersTest(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);
	`)
	if err != nil {
		t.Fatalf("failed to create testing table :: %v", err)
	}
	defer dropTestTable(db, "insertUsersTest")

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	usedID, usedName := uuid.NewString(), "Ashton Babe"

	res, err := data.Insert(ctx, db, "insertUsersTest", data.SQLInsertArgs{
		What:   []string{"id", "name"},
		Values: []interface{}{usedID, usedName},
	})
	if err != nil {
		t.Fatalf("insertion failed res=%v :: err=%v", res, err)
	}

	var (
		id   string
		name string
	)
	err = db.QueryRow("SELECT id, name FROM insertUsersTest WHERE id = ?", usedID).Scan(&id, &name)
	if err == sql.ErrNoRows {
		t.Fatalf("there are no rows with [ID : %s]", usedID)
	} else if err != nil {
		t.Fatal(err)
	}

	if usedID != id && usedName != name {
		t.Fatalf("data mismatch : expected [ID : %s] and [Name : %s], got [ID : %s] and [Name : %s]", usedID, usedName, id, name)
	}

	t.Logf("found [ID : %s] and [Name : %s]", id, name)
}

func TestUpdate(t *testing.T) {
	db, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS updateUsersTest(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);
	`)
	if err != nil {
		t.Fatalf("failed to create testing table :: %v", err)
	}
	defer dropTestTable(db, "updateUsersTest")

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	usedID, usedName := uuid.NewString(), "Ashton Is A Babe"
	updateName := "Ashton Is Not A Babe"

	_, err = db.Exec("INSERT INTO updateUsersTest(id, name) VALUES(?, ?);", usedID, usedName)
	if err != nil {
		t.Fatalf("failed to insert data :: %v", err)
	}

	res, err := data.Update(ctx, db, "updateUsersTest", data.SQLUpdateArgs{
		Set: map[string]interface{}{"name": updateName},
		Where: data.SQLWhereClause{
			Condition: map[string]interface{}{"id": usedID, "name": usedName},
			Operator: "AND",
		},
	})
	if err != nil {
		t.Fatalf("update failed res=%v :: err=%v", res, err)
	}

	var (
		id   string
		name string
	)
	err = db.QueryRow("SELECT id, name FROM updateUsersTest WHERE id = ?", usedID).Scan(&id, &name)
	if err == sql.ErrNoRows {
		t.Fatalf("there are no rows with [ID : %s]", usedID)
	} else if err != nil {
		t.Fatal(err)
	}

	if name != updateName {
		t.Fatalf("data mismatch : expected (updated) [Name : %s], got [Name: %s]", updateName, name)
	}

	t.Logf("found [ID : %s] and [Name : %s]", id, name)
}

func TestDelete(t *testing.T) {
	db, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS deleteUsersTest(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);
	`)
	if err != nil {
		t.Fatalf("failed to create testing table :: %v", err)
	}
	defer dropTestTable(db, "deleteUsersTest")

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	usedID, usedName := uuid.NewString(), "Ashton Is Such A Babe"

	_, err = db.Exec("INSERT INTO deleteUsersTest(id, name) VALUES(?, ?);", usedID, usedName)
	if err != nil {
		t.Fatalf("failed to insert data :: %v", err)
	}
	
	res, err := data.Delete(ctx, db, "deleteUsersTest", data.SQLDeleteArgs{
		Where: data.SQLWhereClause{
			Condition: map[string]interface{}{"id": usedID, "name": usedName},
			Operator: "AND",
		},
	})
	if err != nil {
		t.Fatalf("delete failed res=%v :: err=%v", res, err)
	}

	var (
		id   string
		name string
	)
	err = db.QueryRow("SELECT id, name FROM deleteUsersTest WHERE id = ?", usedID).Scan(&id, &name)
	if err != sql.ErrNoRows {
		t.Fatalf("Found {[ID : %s] and [Name : %s]} - expected nothing", id, name)
	}
}

func TestSelectMany(t *testing.T) {
	db, err := createTestDB()
	testTableName := "selectManyTest"

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL
	);
	`, testTableName))
	if err != nil {
		t.Fatalf("failed to create testing table :: %v", err)
	}
	defer dropTestTable(db, testTableName)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	usedID, usedName := uuid.NewString(), "Ashton Is Such A Babe"

	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s(id, name) VALUES(?, ?);", testTableName), usedID, usedName)
	if err != nil {
		t.Fatalf("failed to insert data :: %v", err)
	}
	
	rows, err := data.SelectMany(ctx, db, testTableName, data.SQLSelectArgs{
		What: []string{"id", "name"},
		Where: data.SQLWhereClause{
			Condition: map[string]interface{}{"id":usedID, "name": usedName},
			Operator: "AND",
		},
		Limit: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !rows.Next() {
		t.Fatal("No rows obtained")
	}
	
	for rows.Next() {
		var (
			id string 
			name string 
		)
		rows.Scan(&id, &name)

		if id != usedID && name != usedName {
			t.Fatalf("Expected [ID : %s] and [Name : %s], but got [ID : %s] and [Name : %s]", usedID, usedName, id, name)
		}
	}
}