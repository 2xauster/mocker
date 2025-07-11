package data

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

func PrepareCreateTableStmt(name string, fields map[string]SQLField) string {
	stmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (`, name)
		
	fieldsSlice := make([]SQLField, 0, len(fields))
	for _, v := range fields {
		fieldsSlice = append(fieldsSlice, v)
	}

	for i := range fieldsSlice {
		v := fieldsSlice[i]
		constraints := v.Constraints

		stmt += fmt.Sprintf(`%s %s %s`, v.Name, v.Datatype, constraints)

		if i < len(fieldsSlice) - 1 {
			stmt += ","
		}
		stmt += "\n"
	}
	stmt += ");"
	return stmt
}

/*
* Generic function to create SQL tables
 */
func CreateTable[T interface{}](ctx context.Context, db *sql.DB, entity T) (string, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", err 
	}

	typ := reflect.TypeOf(entity)
	name := typ.Name()

	fields := ExtractFields(entity)
	if len(fields) <= 0 {
		return "", fmt.Errorf("[func CreateTable] entity of generic type %T has no fields", entity)
	}
	
	stmt := PrepareCreateTableStmt(name, fields)
	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return stmt, fmt.Errorf("[func CreateTable] transaction failed :: %w", err) 
	}

	err = tx.Commit()
	if err != nil {
		return stmt, fmt.Errorf("[func CreateTable] failed to commit :: %w", err)
	}
	return stmt, nil
}