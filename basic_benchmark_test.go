package goany

import (
	"testing"
	"time"
)

func BenchmarkToInt64E(b *testing.B) {
	benchmarks := []struct {
		name  string
		input interface{}
	}{
		{"Nil", nil},
		{"String_Small", "123"},
		{"String_Large", "9223372036854775807"},
		{"Int", 123},
		{"Uint", uint(123)},
		{"Uint64_Large", uint64(9223372036854775807)},
		{"Uint64_Overflow", uint64(9223372036854775808)},
		{"Float32", float32(3.14)},
		{"Float64", 3.14},
		{"Bool_True", true},
		{"Bool_False", false},
		{"Time", time.Now()},
		{"[]uint8_Empty", []uint8{}},
		{"[]uint8_Valid", []uint8("123")},
		{"[]uint8_Invalid", []uint8("abc")},
		{"Struct", struct{}{}},
		{"Slice", []int{1, 2, 3}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = toInt64E(bm.input)
			}
		})
	}
}

func BenchmarkToUint64E(b *testing.B) {
	benchmarks := []struct {
		name  string
		input interface{}
	}{
		{"Nil", nil},
		{"String_Positive", "123"},
		{"String_Negative", "-123"},
		{"String_Large", "18446744073709551615"},
		{"Int_Positive", 123},
		{"Int_Negative", -123},
		{"Uint", uint(123)},
		{"Float32_Positive", float32(3.14)},
		{"Float32_Negative", float32(-3.14)},
		{"Float64_Positive", 3.14},
		{"Float64_Negative", -3.14},
		{"Bool", true},
		{"[]uint8_Empty", []uint8{}},
		{"[]uint8_Valid", []uint8("123")},
		{"[]uint8_Invalid", []uint8("abc")},
		{"Struct", struct{}{}},
		{"Slice", []int{1, 2, 3}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = toUint64E(bm.input)
			}
		})
	}
}

func BenchmarkToFloat64E(b *testing.B) {
	testCases := []struct {
		name  string
		input interface{}
	}{
		{"Float64", 123.45},
		{"Int", 123},
		{"String", "123.45"},
		{"Bool", true},
		{"Slice", []uint8("123")},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = toFloat64E(tc.input)
			}
		})
	}
}

func BenchmarkToStringE(b *testing.B) {
	op := Options{
		timeFormat: "2006-01-02 15:04:05",
	}

	tests := []struct {
		name string
		v    interface{}
		op   Options
	}{
		{"String", "test", op},
		{"Int", 123, op},
		{"Uint", uint(123), op},
		{"Float32", float32(123.45), op},
		{"Float64", float64(123.45), op},
		{"Bool", true, op},
		{"Time", time.Now(), op},
		{"Struct", struct{ Name string }{"test"}, op},
		{"Map", map[string]interface{}{"key": "value"}, op},
		{"Slice", []int{1, 2, 3}, op},
		{"Array", [3]int{1, 2, 3}, op},
		{"ByteSlice", []byte("test"), op},
		{"Nil", nil, op},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := toStringE(tt.v, tt.op)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}

func BenchmarkToBoolE(b *testing.B) {
	benchmarks := []struct {
		name  string
		input interface{}
	}{
		{"Nil", nil},
		{"StringTrue", "true"},
		{"StringFalse", "false"},
		{"StringInvalid", "invalid"},
		{"Int", 1},
		{"Uint", uint(1)},
		{"Float32", float32(1.0)},
		{"Float64", float64(1.0)},
		{"Bool", true},
		{"ByteSliceTrue", []uint8("true")},
		{"ByteSliceFalse", []uint8("false")},
		{"SliceInt", []int{1}},
		{"Struct", struct{}{}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = toBoolE(bm.input)
			}
		})
	}
}
