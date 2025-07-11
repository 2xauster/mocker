package data

import (
	"reflect"
)

// SQLField represents a field in a sql table.
type SQLField struct {
	Name string
	Datatype string
	Constraints string  
}

/*
Function to generate SQL table fields from a golang structure.
It returns a map of [type string : "name of field"]: [type SQLField] 
*/
func ExtractFields[T interface{}](entity T) map[string]SQLField {
	typ := reflect.TypeOf(entity)

	fields := make(map[string]SQLField, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		constraintsTag := field.Tag.Get("cnstr")
		datatypeTag := field.Tag.Get("type")

		fields[field.Name] = SQLField{
			Name: field.Name,
			Datatype: datatypeTag,
			Constraints: constraintsTag,
		}
	}
	return fields
}
