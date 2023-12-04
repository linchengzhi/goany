# [English](https://github.com/linchengzhi/goany/REREAD.MD) | [中文](https://github.com/linchengzhi/goany/README_ZH.MD)
# GoAny: Go 语言的强大类型转换库

`GoAny` 是一个多功能的 Go 语言库，提供了一种方式将不同的数据类型（如基础类型、列表、映射、结构体等）转换为指定的目标类型。这个工具使用反射动态处理类型转换，使其成为各种 Go 应用程序的灵活工具。

## 安装
要安装 GoAny，请使用 `go get` 命令：

```bash
go get github.com/linchengzhi/goany
```
## 示例
GoAny 提供了通用的转换函数 `ToAny(in interface{}, out interface{}, options... Options) error`，能够转换基础类型、列表、映射、结构体等。它还提供了一系列基本类型转换函数，如 `ToInt(v interface{}) int`, `ToString(v interface{}) string`, `ToTimeE(v interface{}, op ...Options) (time.Time, error)`等。以 'E' 结尾的函数会返回结果和错误。
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
PS：更多示例请查看测试。

## 支持的类型
GoAny 支持以下转换：
- #### 基础类型（整数、无符号整数、浮点数、布尔、字符串）
  提供基础类型之间的常见转换，例如： 基础类型有直接的转换函数。
  ```go
  v := goany.ToInt("1") //"1"
  v, err := goany.ToInt64E(123) //123 nil Functions ending with 'E' return an error along with the result.
  v := goany.ToString(1) //"1"
  ```
  PS：如果结构体（除了 time.Time）、列表、映射转为字符串，将会被转换为 JSON。如果列表是 []byte，它将被转换为字符串。
- #### 切片和数组
  提供列表、映射、字符串到列表的转换，字符串必须是 JSON 格式
  ```go
  var in = []interface{}{1, 2}
  var out = []string{}
  err := goany.ToAny(in, &out) //[]string{"1", "2"}
  ```
- #### 映射
  提供列表、映射、结构体、字符串到映射的转换，字符串必须是 JSON 格式
  ```go
  type player struct {
      Id   int    `json:"id"`
      Name string `json:"name"`
  }
  var in = player{Id: 1, Name: "a"}
  var out = make(map[string]interface{})
  err := goany.ToAny(in, &out) //map[string]interface{}{"id": 1, "name": "a"}
  ```
- #### 结构体（包括 time.Time）
  提供列表、映射、结构体、字符串到结构体的转换，字符串必须是 JSON 格式
  ```go
  type player struct {
      Id   int    `json:"id"`
      Name string `json:"name"`
  }
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out = player{}
  err := goany.ToAny(in, &out) //player{Id: 1, Name: "a"}
  ```
- #### 接口
  任何类型都可以转换为接口。如果输入是结构体，并且选项 structToMapDetail 为真，它将被转换为 map[string]interface{}
  ```go
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out interface{}
  err := goany.ToAny(in, &out) //map[string]interface{}{"id": 1, "name": "a"}
  ```

## 选项
- #### location
  时区默认为 "UTC"。
  ```go
  locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
  var in = "2020-10-01 21:06:11"
  op := goany.NewOptions().SetLocation(locationShanghai)
  out, err := goany.ToTimeE(in, *op)
  fmt.Println(out, err) //2020-10-01 21:06:11 +0800 CST, nil
  ```
- #### timeFormat
  时间格式默认为 "2006-01-02 15:04:05"
  ```go
  locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
  in := time.Date(2020, 10, 1, 21, 6, 11, 0, locationShanghai)
  op := goany.NewOptions().SetLocation(locationShanghai).SetTimeFormat(time.RFC3339)
  out, err := goany.ToStringE(in, *op)
  fmt.Println(out, err)//2020-10-01T21:06:11+08:00 <nil>
  ```
- #### mapKeyField
  将列表结构体转换为映射结构体时，使用 mapKeyField 选项指定作为映射键的结构体字段。
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
  映射转换为列表时，默认转换为值列表。当 mapKeyToList 为真时，使用映射键。
  ```go
  var in = map[string]interface{}{"id": 1, "name": "a"}
  var out = make([]interface{}, 0)
  op := goany.NewOptions().SetMapKeyToList(true)
  err := goany.ToAny(in, &out, *op) 
  fmt.Println(out, err)//[id name] nil
  ```
- #### tagName
  指定结构体的标签作为字段名。默认标签是 json。如果没有标签，则使用字段名。
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
  如果 exportedUnExported 值为真，则结构体中不可导出的导出字段可以被导出。如果in包含要导出不可导出的结构字段，则in必须用使用指针
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
  当您想要将结构体深度转换为映射时，将 structToMapDetail 设置为真。如果in包含要导出不可导出的结构字段，则in必须用使用指针
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
  当您想要将一种结构体转换为另一种结构体时，可以使用 assignKey 来指定字段名。
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
  当您想要自定义解析时，可以使用钩子。钩子函数: `func(in interface{}, out reflect.Value) (int, error)`。in 是输入值，out 是输出值。返回值是一个整数，表示解析的状态。
  如果返回 `goany.DecodeContinue`，则继续解析。如果返回 `goany.DecodeSkip`，则跳过当前值。如果返回 `goany.DecodeStop`，停止所有解析。error不为空时，亦停止所有解析。
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
欢迎贡献以改进 ToAny！请随时向仓库提交问题和拉取请求。

## License
GoAny 采用 MIT 许可证授权。详细信息请查看 [LICENSE](https://github.com/linchengzhi/goany/LICENSE) 文件。