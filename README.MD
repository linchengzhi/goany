# [English](https://github.com/linchengzhi/goany/REREAD.MD) | [中文](https://github.com/linchengzhi/goany/README_ZH.MD)

# GoAny: A Powerful Type Conversion Library for Go

`GoAny` is a versatile Go lib that provides a way to convert different data types such as base types, lists, maps, structures, and more, into a specified target type. This utility uses reflection to dynamically handle type conversion, making it a flexible tool for various Go applications.

## Installation

To install GoAny, use the `go get` command:

```bash
go get github.com/linchengzhi/goany
```
## Example
GoAny provides the generic conversion function `ToAny(in interface{}, out interface{}, options... Options) error`, which can convert base types, list, map, struct, etc. It also provides a set of basic type conversion functions, such as `ToInt(v interface{}) int`, `ToString(v interface{}) string`, `ToTimeE(v interface{}, op ...Options) (time.Time, error)`, etc. Functions ending with 'E' return an error along with the result.
```go
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
```
PS:More examples please see test.

## Supported Types
GoAny supports the following conversions:
- #### Basic types (integers, unsigned integers, floats, bool, string)
  Provides common conversions between base types, for example：
  Base types have direct conversion functions.
  ```go
  v := goany.ToInt("1") //"1"
  v, err := goany.ToInt64E(123) //123 nil Functions ending with 'E' return an error along with the result.
  v := goany.ToString(1) //"1"
  ```
  PS:If struct(except time.Time),list,map to string, it will be converted to json. if list is []byte, it will be converted to string
- #### Slices and arrays
  Provides list, map, string to list conversion, string must be json format
  ```go
  var in = []interface{}{1, 2}
  var out = []string{}
  err := goany.ToAny(in, &out) //[]string{"1", "2"}
  ```
- #### Maps
  Provides list, map, struct, string to map conversion, string must be json format
  ```go
  type player struct {
      Id   int    `json:"id"`
      Name string `json:"name"`
  }
  var in = player{Id: 1, Name: "a"}
  var out = make(map[string]interface{})
  err := goany.ToAny(in, &out) //map[string]interface{}{"id": 1, "name": "a"}
  ```
- #### Structs (including time.Time)
  Provides list, map, struct, string to struct conversion, string must be json format
  ```go
  type player struct {
      Id   int    `json:"id"`
      Name string `json:"name"`
  }
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out = player{}
  err := goany.ToAny(in, &out) //player{Id: 1, Name: "a"}
  ```
- #### Interface
  Everything can to interface.If in is struct, option structToMapDetail is true, it will be converted to map[string]interface{}
  ```go
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out interface{}
  err := goany.ToAny(in, &out) //map[string]interface{}{"id": 1, "name": "a"}
  ```
## Options
- #### location
  Time zone default is "UTC".
  ```go
  locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
  var in = "2020-10-01 21:06:11"
  op := goany.NewOptions().SetLocation(locationShanghai)
  out, err := goany.ToTimeE(in, *op)
  fmt.Println(out, err) //2020-10-01 21:06:11 +0800 CST, nil
  ```
- #### timeFormat
  Time format default is "2006-01-02 15:04:05"
  ```go
  locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
  in := time.Date(2020, 10, 1, 21, 6, 11, 0, locationShanghai)
  op := goany.NewOptions().SetLocation(locationShanghai).SetTimeFormat(time.RFC3339)
  out, err := goany.ToStringE(in, *op)
  fmt.Println(out, err)//2020-10-01T21:06:11+08:00 <nil>
  ```
- #### mapKeyField
  When list struct to map struct, use the mapKeyField option to specify which struct field to use as the map key when converting a list of structs to a map.
  ```go
  type player struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
  }
  var in = []player{{Id: 1, Name: "a"}, {Id: 2, Name: "b"}}
  var out = make(map[string]player)
  op := goany.NewOptions().SetMapKeyField("name")
  err := goany.ToAny(in, &out, *op) 
  fmt.Println(out, err)//map[a:{1 a} b:{2 b}] nil
  ```
- #### mapKeyToList
  When map is converted to list, value is converted to list by default. When mapKeyToList is true, map key is used
  ```go
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out = make([]interface{}, 0)
  op := goany.NewOptions().SetMapKeyToList(true)
  err := goany.ToAny(in, &out, *op) 
  fmt.Println(out, err)//[id name] nil
  ```
- #### tagName
  Specify the tag of the struct as the field name. The default tag is json. If there is no tag, use the field name
  ```go
  type player struct {
    Id   int    `bson:"id"`
    Name string `bson:"name"`
  }
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out = player{}
  op := goany.NewOptions().SetTagName("bson")
  err := goany.ToAny(in, &out, *op) 
  fmt.Println(out, err)//player{Id: 1, Name: "a"}
  ```
- #### exportedUnExported
  If the exportedUnExported value is true, exported fields that are not exportable in the structure can be exported.
  If in contains non-exportable structure fields to be exported, in must be a pointer
  ```go
  type player struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
    age  int    `json:"age"` //unexported
  }
  var in = map[string]interface{}{"id": 1, "name": "a", "age": 20}
  var out = player{}
  op := goany.NewOptions().SetExportedUnExported(true)
  err := goany.ToAny(in, &out, *op) 
  fmt.Println(out, err)//player{Id: 1, Name: "a", age: 20}
  ```
- #### structToMapDetail
  When you want to convert structure depth to map, set structToMapDetail to true.If in contains non-exportable structure fields to be exported, in must be used with a use pointer
  ```go
  type player struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
  }
  type account struct {
    Name   string `json:"name"`
    player player `json:"player"` //nest struct
  }
  var in = account{Name: "a", player: player{Id: 1, Name: "a"}}
  var out = make(map[string]interface{})
  op := goany.NewOptions().SetStructToMapDetail(true).SetExportedUnExported(true)
  err := goany.ToAny(&in, &out, *op) // in is pointer
  fmt.Println(out, err)//map[string]interface{}{"name": "a", "player": map[string]interface{}{"id": 1, "name": "a"}}
  ```
- #### assignKey
  When you want to convert structure to another structure, you can use `assignKey` to assign the field name
  ```go
  type player struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
  }
  type student struct {
    Sid  int    `json:"sid"` //id to sid
    Name string `json:"name"`
  }
  in := &player{Id: 1, Name: "a"}
  out := &student{}
  op := goany.NewOptions().SetAssignKey(map[string]string{"id": "sid"})
  err := goany.ToAny(in, out, *op) 
  fmt.Println(out, err)//&student{Sid: 1, Name: "a"}
  ```
- #### hooks
  When you need to customize the parsing process, you can use hooks. The hook function is defined as: `func(in interface{}, out reflect.Value) (int, error)`. Here, in is the input value and out is the output value. The returned integer indicates the parsing state.
  If the hook returns `goany.DecodeContinue`, the parsing continues. If it returns `goany.DecodeSkip`, the current value is skipped. If it returns `goany.DecodeStop`, it halts all parsing. Parsing also stops if the error is not nil.
  ```go
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
  ```
## Contributing
Contributions to improve ToAny are welcome! Please feel free to submit issues and pull requests to the repository.

## License
GoAny is licensed under the MIT license. See the [LICENSE](https://github.com/linchengzhi/goany/LICENSE) file for details.