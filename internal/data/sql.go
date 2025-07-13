package data

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type SQLWhereClause struct {
	Condition map[string]interface{}
	Operator  string // AND, OR (mandatory)
}

/*
* Data to perform "INSERT" operation in SQL.
 */
type SQLInsertArgs struct {
	What   []string
	Values []interface{}
}

// Data to perform "UPDATE" operation in SQL.
type SQLUpdateArgs struct {
	Set   map[string]interface{}
	Where SQLWhereClause
}

// Data to perform "DELETE" operation in SQL.
type SQLDeleteArgs struct {
	Where SQLWhereClause 
}

// Data to perform "SELECT" operation in SQL.
type SQLSelectArgs struct {
	What []string // e.g: id, name, email, etc.
	Where SQLWhereClause
	Limit int // limit the amount of rows to fetch
}

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

		if i < len(fieldsSlice)-1 {
			stmt += ","
		}
		stmt += "\n"
	}
	stmt += ");"
	return stmt
}

func PrepareWhereClause(clause SQLWhereClause) (string, []interface{}, error) {
	values := make([]interface{}, 0, len(clause.Condition))

	whereParts := make([]string, 0, len(clause.Condition))
	for k, v := range clause.Condition {
		whereParts = append(whereParts, fmt.Sprintf(`%s=?`, k))
		values = append(values, v)
	}
	op := strings.ToUpper(clause.Operator)
	if op != "OR" && op != "AND" {
		return "", []interface{}{}, fmt.Errorf("data.SQLWhereClause:Where:Operator should be AND / OR, not %s", op)
	}
	where := "WHERE " + strings.Join(whereParts, " "+op+" ")

	return where, values, nil
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

func Insert(ctx context.Context, db *sql.DB, tableName string, args SQLInsertArgs) (sql.Result, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Insert] failed to begin transaction :: %w", err)
	}
	defer tx.Rollback()

	what := strings.Join(args.What, ", ")
	placeholders := make([]string, len(args.Values))
	for i := range args.Values {
		placeholders[i] = "?"
	}

	stmt := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s);`, tableName, what, strings.Join(placeholders, ", "))

	res, err := tx.ExecContext(ctx, stmt, args.Values...)

	if err != nil {
		return res, fmt.Errorf("[pkg data : func Insert] execution failed :: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return  res, fmt.Errorf("[pkg data : func Insert] failed to commit :: %w", err)
	}
	return res, err
}

func Update(ctx context.Context, db *sql.DB, tableName string, args SQLUpdateArgs) (sql.Result, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Update] failed to begin transaction :: %w", err)
	}
	defer tx.Rollback()

	setParts := make([]string, 0, len(args.Set))
	values := make([]interface{}, 0, len(args.Set))

	for k, v := range args.Set {
		setParts = append(setParts, fmt.Sprintf(`%s=?`, k))
		values = append(values, v)
	}

	where, whereVals, err := PrepareWhereClause(args.Where)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Update] where clause preparation failed :: %w", err)
	}
	values = append(values, whereVals...)

	stmt := fmt.Sprintf(`UPDATE %s SET %s %s`, tableName, strings.Join(setParts, ", "), where)

	res, err := tx.ExecContext(ctx, stmt, values...)
	if err != nil {
		return res, fmt.Errorf("[pkg data : func Update] execution failed :: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return res, fmt.Errorf("[pkg data : func Update] commit failed :: %w", err)
	}

	return res, err
}

func Delete(ctx context.Context, db *sql.DB, tableName string, args SQLDeleteArgs) (sql.Result, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Delete] failed to begin transaction :: %w", err)
	}
	defer tx.Rollback()
	
	where, values, err := PrepareWhereClause(args.Where)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Delete] preparation of where clause failed :: %w", err)
	}

	stmt := fmt.Sprintf(`DELETE FROM %s %s`, tableName, where)
	res, err := tx.ExecContext(ctx, stmt, values...)
	if err != nil {
		return res, fmt.Errorf("[pkg data : func Delete] execution failed :: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return res, fmt.Errorf("[pkg data : func Delete] commit failed :: %w", err)
	}

	return res, err
}

func SelectMany(ctx context.Context, db *sql.DB, tableName string, args SQLSelectArgs) (*sql.Rows, error) {
	what := strings.Join(args.What, ", ")
	where, whereValues, err := PrepareWhereClause(args.Where)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Select] failed to prepare where clause :: %w", err)
	}
	
	stmt := fmt.Sprintf(`SELECT %s FROM %s %s LIMIT %d;`, what, tableName, where, args.Limit)
	rows, err := db.QueryContext(ctx, stmt, whereValues...)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Select] query failed :: %w", err)
	}
	return rows, nil
}