package data

import (
	"reflect"
	"strings"
)

// SQLField represents a field in a SQL table.
type SQLField struct {
	Name        string
	Datatype    string
	Constraints string
	Reference   string // e.g., "Users(ID)" or "Authors(id)"
}

type EntityMetadata struct {
	Name   string
	Fields []SQLField
}

/*
ExtractFields generates SQL table fields from a Go structure.
It returns a map from field name to SQLField definition.
*/
func ExtractFields[T any](entity T, zeroValidation bool) []SQLField {
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	fields := make([]SQLField, 0, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Ensure field values are not zeroes (if required)
		if zeroValidation && val.Field(i).IsZero() {
			continue
		}
		constraintsTag := field.Tag.Get("cnstr")
		datatypeTag := field.Tag.Get("type")
		referenceTag := field.Tag.Get("ref")

		fields = append(fields, SQLField{
			Name:        strings.ToLower(field.Name),
			Datatype:    datatypeTag,
			Constraints: constraintsTag,
			Reference:   referenceTag,
		})
	}
	return fields
}

// Utility function to extract the data of go structs.
func ExtractMeta[T any](entity T, extractRealFields bool) EntityMetadata {
	typ := reflect.TypeOf(entity)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := ExtractFields(entity, extractRealFields)
	name := strings.ToLower(typ.Name())

	return EntityMetadata{
		Name:   name,
		Fields: fields,
	}
}

func ExtractValueSlice[T any](entityData T) []any {
	val := reflect.ValueOf(entityData)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("input must be a struct or pointer to struct")
	}
	result := make([]any, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// Do not append if the value is zero.
		if field.IsZero() {
			continue
		}

		result = append(result, field.Interface())
	}
	return result
}
