package data

import (
	"reflect"
)

// SQLField represents a field in a SQL table.
type SQLField struct {
	Name        string
	Datatype    string
	Constraints string
	Reference   string // e.g., "Users(ID)" or "Authors(id)"
}

/*
ExtractFields generates SQL table fields from a Go structure.
It returns a map from field name to SQLField definition.
*/
func ExtractFields[T any](entity T) map[string]SQLField {
	typ := reflect.TypeOf(entity)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := make(map[string]SQLField, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		constraintsTag := field.Tag.Get("cnstr")
		datatypeTag := field.Tag.Get("type")
		referenceTag := field.Tag.Get("ref")

		fields[field.Name] = SQLField{
			Name:        field.Name,
			Datatype:    datatypeTag,
			Constraints: constraintsTag,
			Reference:   referenceTag,
		}
	}

	return fields
}