package goany

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDecodeMap_Map(t *testing.T) {
	tests := []structTest{
		{
			name:     "Test with map",
			input:    map[string]interface{}{"a": "1", "b": "2"},
			output:   make(map[string]string),
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			name:     "Test with nested map",
			input:    map[string]interface{}{"a": map[string]interface{}{"b": "1"}},
			expected: map[string]interface{}{"a": map[string]interface{}{"b": "1"}},
		},
		{
			name:     "Test with two-dimensional array",
			input:    [][]interface{}{{"1", "2"}, {"3", "4"}},
			expected: map[string]interface{}{"0": []interface{}{"1", "2"}, "1": []interface{}{"3", "4"}},
		},
		{
			name:     "Test with struct containing slice",
			input:    struct{ S []string }{S: []string{"1", "2"}},
			expected: map[string]interface{}{"S": []string{"1", "2"}},
		},
		{
			name:     "Test with unsupported type",
			input:    123,
			expected: nil,
			err:      errors.Errorf(ErrInToOut, 123, "map"),
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

func TestDecodeMap_List(t *testing.T) {
	type account struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []structTest{
		{
			name:     "Test with interface slice",
			input:    []interface{}{"1", "2"},
			expected: map[string]interface{}{"0": "1", "1": "2"},
		},
		{
			name:     "Test with two-dimensional slice",
			input:    [][]interface{}{{"1", "2"}, {"3", "4"}},
			expected: map[string]interface{}{"0": []interface{}{"1", "2"}, "1": []interface{}{"3", "4"}},
		},
		{
			name:     "Test with two-dimensional slice",
			input:    [2][2]interface{}{{"1", "2"}, {"3", "4"}},
			expected: map[string]interface{}{"0": [2]interface{}{"1", "2"}, "1": [2]interface{}{"3", "4"}},
		},
		{
			name:     "Test with struct containing slice",
			input:    struct{ S []string }{S: []string{"1", "2"}},
			expected: map[string]interface{}{"S": []string{"1", "2"}},
		},
		{
			name:     "Test with struct use field to map",
			input:    []*account{{Id: 1, Name: "abc"}, {Id: 2, Name: "def"}},
			op:       NewOptions().SetMapKeyField("id"),
			expected: map[string]interface{}{"1": &account{Id: 1, Name: "abc"}, "2": &account{Id: 2, Name: "def"}},
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

func TestDecodeMap_Struct(t *testing.T) {
	type player struct {
		Id   int    `gorm:"column:gorm_id" json:"id"`
		Name string `bson:"bson_name" json:"name"`
	}

	type AnonymousData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type anonymousNoExport struct {
		Name string `json:"name"`
	}

	type account struct {
		Name   string `json:"name"`
		Player player `json:"player"` //nest struct
	}

	type accountWithAnonymous struct {
		Name   string  `json:"name"`
		Player *player `json:"player"` //nest struct
		AnonymousData
	}

	type accountWithUnEx struct {
		name string `json:"name"`
		anonymousNoExport
	}

	var Vtime = time.Now()
	var Rint = func() *int { i := 1; return &i }()

	tests := []structTest{
		{
			name:  "Test with basic struct",
			input: BasicType{Vint: 1, Vstr: "a", Vbool: true, Vfloat: 1.12, Vtime: Vtime, Rint: Rint, Nstr: "b"},
			expected: map[string]interface{}{"vint": 1, "vstr": "a", "vbool": true, "vfloat": 1.12, "vtime": Vtime,
				"rint": Rint, "nstr": json.Number("b")},
		},
		{
			name:   "Test with basic struct to string, rint set nil",
			input:  BasicType{Vint: 1, Vstr: "a", Vbool: true, Vfloat: 1.12, Vtime: Vtime, Nstr: "b"},
			output: make(map[string]string),
			expected: map[string]string{"vint": "1", "vstr": "a", "vbool": "true", "vfloat": "1.12", "rint": "",
				"vtime": Vtime.Format("2006-01-02 15:04:05"), "nstr": "b"},
		},
		{
			name:     "Test with basic type 2, unexported field, ignore field, no json tag",
			input:    &BasicType2{unexported: "unexported", Ignore: "ignore", NoTag: true},
			expected: map[string]interface{}{"NoTag": true},
		},
		{
			name:     "Test with basic type 2, export unexported field",
			input:    &BasicType2{unexported: "unexported", Ignore: "ignore", NoTag: true},
			op:       NewOptions().SetExportedUnExported(true),
			expected: map[string]interface{}{"unexported": "unexported", "NoTag": true},
		},
		{
			name:     "Test with basic type 2, export unexported field, to string",
			input:    &BasicType2{unexported: "unexported", Ignore: "ignore", NoTag: true},
			output:   make(map[string]string),
			op:       NewOptions().SetExportedUnExported(true),
			expected: map[string]string{"unexported": "unexported", "NoTag": "true"},
		},
		{
			name:     "Test with struct",
			input:    account{Name: "a", Player: player{Name: "b"}},
			expected: map[string]interface{}{"name": "a", "player": player{Id: 0, Name: "b"}},
		},
		{
			name:     "Test with struct detail",
			input:    account{Name: "a", Player: player{Name: "b"}},
			op:       NewOptions().SetStructToMapDetail(true),
			expected: map[string]interface{}{"name": "a", "player": map[string]interface{}{"id": 0, "name": "b"}},
		},
		{
			name:     "Test with struct detail, export anonymous",
			input:    accountWithAnonymous{Name: "a", AnonymousData: AnonymousData{Name: "b", Age: 1}},
			op:       NewOptions().SetStructToMapDetail(true).SetExportedUnExported(true),
			expected: map[string]interface{}{"name": "a", "player": nil, "AnonymousData": map[string]interface{}{"name": "b", "age": 1}},
		},
		{
			name:  "Test with struct export unexported, input no ptr",
			input: accountWithUnEx{name: "a"},
			op:    NewOptions().SetStructToMapDetail(true).SetExportedUnExported(true),
			err:   ErrInNotPtr,
		},
		{
			name:     "Test with struct export unexported, input is ptr",
			input:    &accountWithUnEx{name: "a", anonymousNoExport: anonymousNoExport{Name: "b"}},
			op:       NewOptions().SetStructToMapDetail(true).SetExportedUnExported(true),
			expected: map[string]interface{}{"name": "a", "anonymousNoExport": map[string]interface{}{"name": "b"}},
		},

		{
			name:     "Test with struct, set gorm as tag",
			input:    player{Id: 1, Name: "a"},
			op:       NewOptions().SetTagName("gorm"),
			expected: map[string]interface{}{"gorm_id": 1, "Name": "a"},
		},
		{
			name:     "Test with struct, set bson as tag",
			input:    player{Id: 1, Name: "a"},
			op:       NewOptions().SetTagName("bson"),
			expected: map[string]interface{}{"Id": 1, "bson_name": "a"},
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

func TestDecodeMap_String(t *testing.T) {
	tests := []structTest{
		{
			name:     "Test with string map",
			input:    `{"a": "1", "b": "2"}`,
			expected: map[string]interface{}{"a": "1", "b": "2"},
		},
		{
			name:     "Test with string slice",
			input:    `["1", "2"]`,
			expected: map[string]interface{}{"0": "1", "1": "2"},
		},
		{
			name:     "Test with string slice struct, nest",
			input:    `[{"a": "1"}, {"b": "2", "c": {"d": "3"}}]`,
			expected: map[string]interface{}{"0": map[string]interface{}{"a": "1"}, "1": map[string]interface{}{"b": "2", "c": map[string]interface{}{"d": "3"}}},
		},
		{
			name:     "Test with string slice strut, to slice map",
			input:    `[{"a": "1"}, {"b": "2"}]`,
			output:   make([]map[string]string, 0),
			expected: []map[string]string{{"a": "1"}, {"b": "2"}},
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

func TestDecodeMap_List2(t *testing.T) {
	var ints = []string{"1", "2", "3"}
	var m = make([]map[int]int, 0)
	err := ToAny(ints, &m)
	assert.Error(t, err)
}
