package data_test

import (
	"reflect"
	"testing"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/google/uuid"
)

type TestUserEntity struct {
	ID string `type:"TEXT" cnstr:"PRIMARY KEY"`
	Name string `type:"TEXT" cnstr:"NOT NULL"`
}

func TestExtractFields(t *testing.T) {
	user := TestUserEntity{}
	fields := data.ExtractFields(user, false)
	
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

func TestExtractValueSlice(t *testing.T) {
	usedID, usedName := uuid.NewString(), "Ashton Is A Slice"
	user := TestUserEntity{
		ID: usedID,
		Name: usedName,
	}
	values := data.ExtractValueSlice(user)

	if len(values) < 2 {
		t.Errorf("Expected a slice of length 1, but got %d", len(values))
	}
	
	if values[0] != usedID && values[1] != usedName {
		t.Errorf("Values mismatch, expected [ID : %s] and [Name : %s] but got [any : %s] and [any : %s]", usedID, usedName, values[0], values[1])
	}
}