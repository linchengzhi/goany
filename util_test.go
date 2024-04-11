package goany

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name string `json:"name" gorm:"column:name"`
	Age  int    `json:"age" gorm:"column:age"`
}

type TestStruct2 struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func assertFieldName(t *testing.T, field reflect.StructField, tag, expected string) {
	t.Helper()
	fieldName := GetFieldNameByTag(field, tag)
	if fieldName != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, fieldName)
	}
}

func TestGetFieldNameByTag(t *testing.T) {
	typeOfTestStruct := reflect.TypeOf(TestStruct{})
	typeOfTestStruct2 := reflect.TypeOf(TestStruct2{})

	t.Run("json tag", func(t *testing.T) {
		assertFieldName(t, typeOfTestStruct.Field(0), "json", "name")
		assertFieldName(t, typeOfTestStruct.Field(1), "json", "age")
	})

	t.Run("gorm tag", func(t *testing.T) {
		assertFieldName(t, typeOfTestStruct.Field(0), "gorm", "name")
		assertFieldName(t, typeOfTestStruct.Field(1), "gorm", "age")
	})

	t.Run("bson tag", func(t *testing.T) {
		assertFieldName(t, typeOfTestStruct2.Field(0), "bson", "name")
		assertFieldName(t, typeOfTestStruct2.Field(1), "bson", "age")
	})

	t.Run("no tag", func(t *testing.T) {
		assertFieldName(t, typeOfTestStruct.Field(0), "", "Name")
		assertFieldName(t, typeOfTestStruct.Field(1), "", "Age")
	})
}
