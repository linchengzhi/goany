package goany

import (
	"reflect"
	"testing"
	"time"
)

func BenchmarkToType2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var d *int
		v := Indirect(d)
		reflect.TypeOf(v).Kind()
	}
}

func BenchmarkToAny(b *testing.B) {
	type Account struct {
		Account  string
		Password string `json:"-"`
	}

	type player struct {
		Account
		Id       int
		Name     string
		Age      int
		Nums     []string
		Birthday time.Time
		CreateAt time.Time
	}

	type student struct {
		Account
		Id       int
		Name     string
		Age      int
		Nums     []int
		Birthday string
		CreateAt int64
	}
	var p = player{
		Account:  Account{Account: "account", Password: "password"},
		Id:       1,
		Name:     "name",
		Age:      18,
		Nums:     []string{"1", "2", "3"},
		Birthday: time.Now(),
		CreateAt: time.Now(),
	}

	var s student

	for i := 0; i < b.N; i++ {
		err := ToAny(p, &s)
		if err != nil {
			b.Errorf("ToAny error: %v", err)
		}
	}
}
