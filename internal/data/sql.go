package data

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

type SQLWhereClause struct {
	Where any
}

type SQLEntity interface{}

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
	What  []string // e.g: id, name, email, etc.
	Where SQLWhereClause
	Limit int // limit the amount of rows to fetch
}

func PrepareCreateTableStmt(name string, fields []SQLField) string {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\n", name)

	fieldLines := make([]string, 0, len(fields))
	foreignKeys := make([]string, 0)

	for _, v := range fields {
		line := fmt.Sprintf(`%s %s %s`, v.Name, v.Datatype, v.Constraints)

		if v.Reference != "" {
			foreignKeys = append(foreignKeys,
				fmt.Sprintf(`FOREIGN KEY(%s) REFERENCES %s`, v.Name, v.Reference))
		}

		fieldLines = append(fieldLines, line)
	}

	allLines := append(fieldLines, foreignKeys...)
	stmt += strings.Join(allLines, ",\n")
	stmt += "\n);"

	return stmt
}

func PrepareWhereClause(clause SQLWhereClause) (string, []any, error) {
	val := ExtractFields(clause.Where, true)
	values := ExtractValueSlice(clause.Where)

	whereParts := make([]string, 0, len(val))
	for _, v := range val {
		whereParts = append(whereParts, fmt.Sprintf(`%s=?`, v.Name))
	}

	op := "AND"
	where := "WHERE " + strings.Join(whereParts, " "+op+" ")

	return where, values, nil
}

func PrepareWhat(fields []SQLField) string {
	what := make([]string, len(fields))
	for i, v := range fields {
		name := strings.ToLower(v.Name)
		what[i] = name
	}

	return strings.Join(what, ", ")
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
	name := strings.ToLower(typ.Name())

	fields := ExtractFields(entity, false)
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

func Insert[Entity any](ctx context.Context, db *sql.DB, entityData Entity) (sql.Result, error) {
	var entity Entity
	meta := ExtractMeta(entity, false)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Insert] failed to begin transaction :: %w", err)
	}
	defer tx.Rollback()

	what := PrepareWhat(meta.Fields)
	placeholders := make([]string, len(meta.Fields))
	for i := range meta.Fields {
		placeholders[i] = "?"
	}

	values := ExtractValueSlice(entityData)

	stmt := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s);`, meta.Name, what, strings.Join(placeholders, ", "))

	res, err := tx.ExecContext(ctx, stmt, values...)

	if err != nil {
		return res, fmt.Errorf("[pkg data : func Insert] execution failed :: %w", err)
	}

	err = tx.Commit()
	err = SugarifyErrors(err)
	
	if err != nil {
		return res, fmt.Errorf("[pkg data : func Insert] failed to commit :: %w", err)
	}
	return res, err
}

func Update[Entity any](ctx context.Context, db *sql.DB, entityData Entity, where SQLWhereClause) (sql.Result, error) {
	meta := ExtractMeta(entityData, true)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Update] failed to begin transaction :: %w", err)
	}
	defer tx.Rollback()

	values := ExtractValueSlice(entityData)
	setParts := make([]string, 0, len(meta.Fields))

	for _, v := range values {
		slog.Info("Value", "val", v)
	}
	for _, v := range meta.Fields {
		setParts = append(setParts, fmt.Sprintf("%s=?", v.Name))
	}

	clause, whereVals, err := PrepareWhereClause(where)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Update] where clause preparation failed :: %w", err)
	}
	values = append(values, whereVals...)

	stmt := fmt.Sprintf(`UPDATE %s SET %s %s`, meta.Name, strings.Join(setParts, ", "), clause)

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

func Delete(ctx context.Context, db *sql.DB, where SQLWhereClause) (sql.Result, error) {	
	meta := ExtractMeta(where.Where, false)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Delete] failed to begin transaction :: %w", err)
	}
	defer tx.Rollback()

	clause, values, err := PrepareWhereClause(where)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Delete] preparation of where clause failed :: %w", err)
	}

	stmt := fmt.Sprintf(`DELETE FROM %s %s`, meta.Name, clause)
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

	// Minimum limit is 10, ofcourse.
	if args.Limit == 0{
		args.Limit = 10
	}

	stmt := fmt.Sprintf(`SELECT %s FROM %s %s LIMIT %d;`, what, tableName, where, args.Limit)
	rows, err := db.QueryContext(ctx, stmt, whereValues...)
	if err != nil {
		return nil, fmt.Errorf("[pkg data : func Select] query failed :: %w", err)
	}
	return rows, nil
}
