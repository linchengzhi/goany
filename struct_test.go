package goany

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type structTest struct {
	name     string
	input    interface{}
	output   interface{}
	expected interface{}
	err      error
	op       *Options
}

type BasicType struct {
	Vint   int         `bson:"vintb" json:"vint"` //bson tag name is vintb，for test
	Vstr   string      `json:"vstr"`
	Vbool  bool        `json:"vbool"`
	Vfloat float64     `json:"vfloat"`
	Vtime  time.Time   `json:"vtime"`
	Rint   *int        `json:"rint"`
	Nstr   json.Number `json:"nstr"`
}

type BasicType2 struct {
	unexported string `json:"unexported"` // unexported field
	Ignore     string `json:"-"`          //ignore
	NoTag      bool   //no tag
}

func ptrInt(i int) *int {
	return &i
}

// helper function to get a pointer to a string
func ptrString(s string) *string {
	return &s
}

func TestDecodeStruct_Map(t *testing.T) {
	type player struct {
		Id   int    `gorm:"column:gorm_id" json:"id"`
		Name string `bson:"bson_name" json:"name"`
	}

	type AnonData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type anonymousNoExport struct {
		Name string `json:"name"`
		age  int    `json:"age"`
	}

	type account struct {
		Name   string `json:"name"`
		Player player `json:"player"` //nest struct
	}

	type accountWithAnonymous struct {
		Name   string  `json:"name"`
		Player *player `json:"player"` //nest struct
		AnonData
	}

	type accountWithUnEx struct {
		name string `json:"name"`
		anonymousNoExport
	}

	tests := []structTest{
		{
			name:     "Test with map",
			input:    map[string]interface{}{"vint": 2, "vstr": "abc", "vbool": true, "vfloat": 1.1, "vtime": "2020-01-01 00:00:00", "rint": 1, "nstr": "1"},
			output:   new(BasicType),
			expected: &BasicType{Vint: 2, Vstr: "abc", Vbool: true, Vfloat: 1.1, Vtime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rint: ptrInt(1), Nstr: "1"},
		},
		{
			name:     "Test with map with default value",
			input:    map[string]interface{}{},
			output:   new(BasicType),
			expected: &BasicType{Vint: 0, Vstr: "", Vbool: false, Vfloat: 0, Vtime: time.Time{}, Rint: nil, Nstr: ""},
		},
		{
			name:     "Test with map, special field",
			input:    map[string]interface{}{"unexported": "a", "ignore": "b", "NoTag": true},
			output:   new(BasicType2),
			expected: &BasicType2{unexported: "", Ignore: "", NoTag: true},
		},
		{
			name:     "Test with map, nest map",
			input:    map[string]interface{}{"name": "a", "player": map[string]interface{}{"id": 1, "name": "b"}},
			output:   new(account),
			expected: &account{Name: "a", Player: player{Name: "b", Id: 1}},
		},
		{
			name:     "Test with map, anon",
			input:    map[string]interface{}{"name": "a", "age": 1},
			output:   new(accountWithAnonymous),
			expected: &accountWithAnonymous{Name: "a", AnonData: AnonData{Age: 1}},
		},
		{
			name:     "Test with map, anon，unexported",
			input:    map[string]interface{}{"name": "a", "age": 1},
			output:   new(accountWithUnEx),
			expected: &accountWithUnEx{},
		},
		{
			name:     "Test with map, anon，export unexported",
			input:    map[string]interface{}{"name": "a", "age": 1},
			output:   new(accountWithUnEx),
			op:       NewOptions().SetExportedUnExported(true),
			expected: &accountWithUnEx{name: "a", anonymousNoExport: anonymousNoExport{age: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = make(map[string]interface{})
			}
			if tt.op == nil {
				tt.op = NewOptions()
			}
			var err error
			err = ToAny(tt.input, &result, *tt.op)
			if tt.err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDecodeStruct_Struct(t *testing.T) {
	type player struct {
		Id   int    `gorm:"column:gorm_id" json:"id"`
		Name string `bson:"bson_name" json:"name"`
	}

	type anonymousData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type anonymousNoExport struct {
		Name string `json:"name"`
		age  int    `json:"age"`
	}

	type account struct {
		Name   string `json:"name"`
		Player player `json:"player"` //nest struct
	}

	type accountWithAnno struct {
		Name   string  `json:"name"`
		Player *player `json:"player"` //nest struct
		anonymousData
	}

	type accountWithUnEx struct {
		name string `json:"name"`
		anonymousNoExport
	}

	type accountNest struct {
		Account account `json:"account"`
	}

	type AnonymousInterface struct {
		name interface{} `json:"name"`
	}

	tests := []structTest{
		{
			name:     "Test with struct",
			input:    &BasicType{Vint: 2, Vstr: "abc", Vbool: true, Vfloat: 1.1, Vtime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rint: ptrInt(1), Nstr: "1"},
			output:   new(BasicType),
			expected: &BasicType{Vint: 2, Vstr: "abc", Vbool: true, Vfloat: 1.1, Vtime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rint: ptrInt(1), Nstr: "1"},
		},
		{
			name:     "Test with struct with default value",
			input:    &BasicType{},
			output:   new(BasicType),
			expected: &BasicType{Vint: 0, Vstr: "", Vbool: false, Vfloat: 0, Vtime: time.Time{}, Rint: nil, Nstr: ""},
		},
		{
			name:     "Test with struct, special field",
			input:    &BasicType2{unexported: "a", Ignore: "b", NoTag: true},
			output:   new(BasicType2),
			expected: &BasicType2{unexported: "", Ignore: "", NoTag: true},
		},
		{
			name:     "Test with struct, special field, export unexported",
			input:    &BasicType2{unexported: "a", Ignore: "b", NoTag: true},
			output:   new(BasicType2),
			op:       NewOptions().SetExportedUnExported(true),
			expected: &BasicType2{unexported: "a", Ignore: "", NoTag: true},
		},
		{
			name:     "Test with struct, nest struct",
			input:    &account{Name: "a", Player: player{Name: "b", Id: 1}},
			output:   new(account),
			expected: &account{Name: "a", Player: player{Name: "b", Id: 1}},
		},
		{
			name:     "Test with struct, no some struct",
			input:    &player{Name: "a", Id: 1},
			output:   new(account),
			expected: &account{Name: "a"},
		},
		{
			name:     "Test with map, anno，unexported",
			input:    map[string]interface{}{"name": "a", "age": 1},
			output:   new(accountWithUnEx),
			expected: &accountWithUnEx{},
		},
		{
			name:     "Test with map, anno，export unexported",
			input:    map[string]interface{}{"name": "a", "age": 1},
			output:   new(accountWithUnEx),
			op:       NewOptions().SetExportedUnExported(true),
			expected: &accountWithUnEx{name: "a", anonymousNoExport: anonymousNoExport{age: 1}},
		},
		{
			name:     "Test with struct,ptr",
			input:    &accountWithAnno{Name: "a", Player: &player{Name: "b", Id: 1}},
			output:   new(account),
			op:       NewOptions().SetExportedUnExported(true),
			expected: &account{Name: "a", Player: player{Name: "b", Id: 1}},
		},
		{
			name:     "Test with struct, more nest",
			input:    &accountNest{Account: account{Name: "a", Player: player{Name: "b", Id: 1}}},
			output:   new(accountNest),
			op:       NewOptions().SetExportedUnExported(true),
			expected: &accountNest{Account: account{Name: "a", Player: player{Name: "b", Id: 1}}},
		},
		{
			name:     "Test with struct, unexported interface",
			input:    &AnonymousInterface{name: "a"},
			output:   new(AnonymousInterface),
			op:       NewOptions().SetExportedUnExported(true),
			expected: &AnonymousInterface{name: "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = make(map[string]interface{})
			}
			if tt.op == nil {
				tt.op = NewOptions()
			}
			var err error
			err = ToAny(tt.input, &result, *tt.op)
			if tt.err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDecodeStruct_String(t *testing.T) {
	type abc struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	type abc2 struct {
		A string `json:"a"`
		B string `json:"b"`
		C abc    `json:"c"`
	}

	tests := []structTest{
		{
			name:     "Test with string map",
			input:    `{"a": "1", "b": "2"}`,
			output:   new(abc),
			expected: &abc{A: "1", B: "2"},
		},
		{
			name:     "Test with string slice struct, nest",
			input:    `{"a": "1", "b": "2", "c": {"a": "3"}}`,
			output:   new(abc2),
			expected: &abc2{A: "1", B: "2", C: abc{A: "3"}},
		},
		{
			name:     "Test with string slice strut, to slice map",
			input:    `[{"a": "1"}, {"b": "2"}]`,
			output:   make([]*abc, 0),
			expected: []*abc{{A: "1"}, {B: "2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = make(map[string]interface{})
			}
			if tt.op == nil {
				tt.op = NewOptions()
			}
			var err error
			err = ToAny(tt.input, &result, *tt.op)
			if tt.err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestHook_1(t *testing.T) {
	type A struct {
		Name string   `json:"name"`
		List []string `json:"list"`
	}

	type B struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	hook := func(in interface{}, out reflect.Value) (int, error) {
		if reflect.TypeOf(in).Kind() == reflect.Slice && out.Kind() == reflect.Int {
			out.SetInt(int64(len(in.([]string))))
			return DecodeSkip, nil
		}
		return DecodeContinue, nil
	}

	assignKey := make(map[string]string)
	assignKey["list"] = "count"

	a := A{Name: "a", List: []string{"a", "b"}}
	b := B{}
	err := ToAny(a, &b, *NewOptions().AddHook(hook).SetAssignKey(assignKey))
	assert.NoError(t, err)
	assert.Equal(t, 2, b.Count)
}

func TestHook_2(t *testing.T) {
	type A struct {
		Name string `json:"name"`
	}

	type B struct {
		Name string `json:"name"`
	}
	hook := func(in interface{}, out reflect.Value) (int, error) {
		inType, inVal := ReflectTypeValue(in)
		if inType.Kind() == reflect.Struct {
			for i := 0; i < inType.NumField(); i++ {
				if inType.Field(i).Name == "Name" {
					inVal.Field(i).SetString(inVal.Field(i).String() + "_test")
				}
			}
		}
		return DecodeContinue, nil
	}

	a := A{Name: "a"}
	b := B{}
	err := ToAny(&a, &b, *NewOptions().AddHook(hook))
	assert.NoError(t, err)
	assert.Equal(t, "a_test", b.Name)
}

func TestDecodeStruct_Anon(t *testing.T) {
	type Account struct {
		id int
	}
	type Player struct {
		Name string
		age  int
	}
	type PlayerAccount struct {
		Account
		Name string
		age  int
	}
	type Anon struct {
		Player
		Num int
		day int
	}
	type Players []*Player
	type AnonWithPtr struct {
		*Account
		*Player
		Num int
		day int
	}
	type AnonWithSlice struct {
		Players
		Num int
		day int
	}
	type AnonWithSlicePtr struct {
		*Players
		Num int
		day int
	}
	type AnonTwoTier struct {
		PlayerAccount
		Num int
		day int
	}

	tests := []structTest{
		{
			name:     "Test with anon",
			input:    &Anon{Num: 1, day: 2, Player: Player{Name: "a", age: 3}},
			output:   new(Anon),
			expected: &Anon{Num: 1, day: 2, Player: Player{Name: "a", age: 3}},
		},
		{
			name:     "Test with anon ptr",
			input:    &AnonWithPtr{Num: 1, day: 2, Player: &Player{Name: "a", age: 3}},
			output:   new(AnonWithPtr),
			expected: &AnonWithPtr{Num: 1, day: 2, Player: &Player{Name: "a", age: 3}},
		},
		{
			name:     "Test with anon slice",
			input:    &AnonWithSlice{Num: 1, day: 2, Players: []*Player{{Name: "a", age: 3}}},
			output:   new(AnonWithSlice),
			expected: &AnonWithSlice{Num: 1, day: 2, Players: []*Player{{Name: "a", age: 3}}},
		},
		{
			name:     "Test with anon slice ptr",
			input:    &AnonWithSlicePtr{Num: 1, day: 2, Players: &Players{{Name: "a", age: 3}}},
			output:   new(AnonWithSlicePtr),
			expected: &AnonWithSlicePtr{Num: 1, day: 2, Players: &Players{{Name: "a", age: 3}}},
		},
		{
			name:     "Test with anon two tier",
			input:    &AnonTwoTier{Num: 1, day: 2, PlayerAccount: PlayerAccount{Name: "a", age: 3, Account: Account{id: 4}}},
			output:   new(AnonTwoTier),
			expected: &AnonTwoTier{Num: 1, day: 2, PlayerAccount: PlayerAccount{Name: "a", age: 3, Account: Account{id: 4}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = make(map[string]interface{})
			}
			if tt.op == nil {
				tt.op = NewOptions().SetExportedUnExported(true)
			}
			var err error
			err = ToAny(tt.input, &result, *tt.op)
			if tt.err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
