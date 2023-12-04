package goany

import (
	"reflect"
	"testing"
)

func TestIndirect(t *testing.T) {
	t.Run("non-nil pointer input", func(t *testing.T) {
		var a *int
		b := 10
		a = &b
		if Indirect(a) != b {
			t.Errorf("Expected %v, but got %v", b, Indirect(a))
		}
	})

	t.Run("nil input", func(t *testing.T) {
		if Indirect(nil) != nil {
			t.Errorf("Expected nil, but got %v", Indirect(nil))
		}
	})
}

func TestReflectTypeValue(t *testing.T) {
	t.Run("non-nil pointer input", func(t *testing.T) {
		var a *int
		b := 10
		a = &b
		rt, rv := ReflectTypeValue(a)
		if rt.Kind() != reflect.Int || rv.Interface() != b {
			t.Errorf("Expected kind %v and value %v, but got kind %v and value %v", reflect.Int, b, rt.Kind(), rv.Interface())
		}
	})

	t.Run("nil input", func(t *testing.T) {
		rt, rv := ReflectTypeValue(nil)
		if rt != nil || rv.IsValid() {
			t.Errorf("Expected nil type and invalid value, but got type %v and value %v", rt, rv)
		}
	})
}

func TestGetFieldNameByTag(t *testing.T) {
	type Test struct {
		Name string `json:"name" gorm:"column:name"`
		Age  int    `json:"age" gorm:"column:age"`
	}
	tt := Test{Name: "Test", Age: 10}

	t.Run("json tag", func(t *testing.T) {
		field, _ := reflect.TypeOf(tt).FieldByName("Name")
		if GetFieldNameByTag(field, "json") != "name" {
			t.Errorf("Expected 'name', but got %v", GetFieldNameByTag(field, "json"))
		}

		field, _ = reflect.TypeOf(tt).FieldByName("Age")
		if GetFieldNameByTag(field, "json") != "age" {
			t.Errorf("Expected 'age', but got %v", GetFieldNameByTag(field, "json"))
		}
	})

	t.Run("gorm tag", func(t *testing.T) {
		field, _ := reflect.TypeOf(tt).FieldByName("Name")
		if GetFieldNameByTag(field, "gorm") != "name" {
			t.Errorf("Expected 'name', but got %v", GetFieldNameByTag(field, "gorm"))
		}

		field, _ = reflect.TypeOf(tt).FieldByName("Age")
		if GetFieldNameByTag(field, "gorm") != "age" {
			t.Errorf("Expected 'age', but got %v", GetFieldNameByTag(field, "gorm"))
		}
	})
}
