package goany

import (
	"reflect"
	"testing"
)

func BenchmarkGetFieldNameByTag(b *testing.B) {
	tt := TestStruct{Name: "Test", Age: 10}
	field, _ := reflect.TypeOf(tt).FieldByName("Name")

	tests := []struct {
		name string
		tag  string
	}{
		{"JSONTag", "json"},
		{"GORMTag", "gorm"},
		{"NoTag", ""},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer() // 重置计时器
			for i := 0; i < b.N; i++ {
				GetFieldNameByTag(field, tt.tag)
			}
		})
	}
}
