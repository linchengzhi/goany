package goany

import (
	"reflect"
	"testing"
)

// 基准测试函数
func BenchmarkDeepInFields(b *testing.B) {
	// 创建一个实例和选项配置
	tt := TestStruct{Name: "Test", Age: 10}
	op := Options{tagName: "json", exportedUnExported: false}

	// 进行基准测试
	for i := 0; i < b.N; i++ {
		_, err := deepInFields(tt, "", op)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeepOutFields(b *testing.B) {
	// 创建一个实例和选项配置
	tt := TestStruct{Name: "Test", Age: 10}
	op := Options{tagName: "json", exportedUnExported: false}
	outVal := reflect.ValueOf(tt)

	// 进行基准测试
	for i := 0; i < b.N; i++ {
		deepOutFields(outVal, "", op)
	}
}
