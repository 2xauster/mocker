package data_test

import (
	"reflect"
	"testing"

	"github.com/ashtonx86/mocker/internal/data"
)

type TestUserEntity struct {
	ID string `type:"TEXT" cnstr:"PRIMARY KEY"`
	Name string `type:"TEXT" cnstr:"NOT NULL"`
}

func TestExtractFields(t *testing.T) {
	user := TestUserEntity{}
	fields := data.ExtractFields(user)
	
	if len(fields) != 2 {
		t.Errorf("Length of TestUserEntity should have been 2, not %d", len(fields))
	}

	typ := reflect.TypeOf(user)

	t.Log("Extracting fields of TestUserEntity")
	t.Log("--------------[START]----------------")

	iter := 0
	for _, v := range fields {
		t.Logf("Name: %s, Dataatype: %s, Constraints: %s", v.Name, v.Datatype, v.Constraints)

		f := typ.Field(iter)

		if f.Name != v.Name {
			t.Errorf("Field name mismatch :: Struct name : %s | Extracted name : %s", f.Name, v.Name)
		}
		iter++
	}

	t.Log("-----------[END]---------------")
}