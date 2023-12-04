package example

import (
	"fmt"
	"github.com/linchengzhi/goany"
	"reflect"
	"testing"
	"time"
)

func TestToAny_Example(t *testing.T) {
	var err error

	// string to int
	vint := goany.ToInt("123")
	fmt.Println(vint) //123

	vint64, err := goany.ToInt64E("123") //An E ending indicates that an error is returned
	fmt.Println(vint64, err)             //123 nil

	//nil to int, if nil, return default value
	vnil := goany.ToInt(nil)
	fmt.Println(vnil) //0

	// string to time, with options
	op := goany.NewOptions().SetLocation(time.UTC)
	vtime := goany.ToTime("2020-10-01 21:06:11", *op)
	fmt.Println(vtime) //2020-10-01 21:06:11 +0000 UTC

	// int to string
	var int1, str1 = 123, ""
	err = goany.ToAny(int1, &str1)
	fmt.Println(str1, err) //123 nil

	// time to string
	var time2, str2 = time.Date(2020, 10, 1, 21, 6, 11, 0, time.UTC), ""
	err = goany.ToAny(time2, &str2)
	fmt.Println(str2, err) //2020-10-01 21:06:11 nil
	str3, err := goany.ToStringE(time2)
	fmt.Println(str3, err) //2020-10-01 21:06:11 nil

	//map to struct
	type Person struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var m1 = map[string]interface{}{
		"id":   1,
		"name": "John",
		"age":  20,
	}
	var p1 Person
	err = goany.ToAny(m1, &p1)
	fmt.Println(p1, err) //{1, John 20} nil
}

func TestToAny_Hook(t *testing.T) {
	type A struct {
		Name string `json:"name"`
	}
	type B struct {
		Name string `json:"name"`
	}
	hook := func(in interface{}, out reflect.Value) (int, error) {
		inType, inVal := goany.ReflectTypeValue(in)
		if inType.Kind() == reflect.Struct {
			for i := 0; i < inType.NumField(); i++ {
				if inType.Field(i).Name == "Name" {
					inVal.Field(i).SetString(inVal.Field(i).String() + "_test")
				}
			}
		}
		return goany.DecodeContinue, nil
	}
	a := A{Name: "a"}
	b := B{}
	err := goany.ToAny(&a, &b, *goany.NewOptions().AddHook(hook))
	fmt.Println(b, err) //{a_test} <nil>
}
