package goany

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToList_Map(t *testing.T) {
	type account struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var tests = []structTest{
		{
			name:     "Test without options, default use value as list",
			input:    map[string]interface{}{"a": "1"},
			expected: []string{"1"},
		},
		{
			name:     "Test with options use key as list",
			input:    map[string]interface{}{"a": "1"},
			op:       NewOptions().SetMapKeyToList(true),
			expected: []string{"a"},
		},
		{
			name:     "Test with map struct as list",
			input:    map[string]interface{}{"a": &account{Id: 1, Name: "abc"}},
			output:   []*account{},
			expected: []*account{{Id: 1, Name: "abc"}},
		},
		{
			name:     "Test struct to slice",
			input:    account{Id: 1, Name: "abc"},
			output:   []string{},
			expected: []string{},
			err:      errors.Errorf(ErrInToOut, account{Id: 1, Name: "abc"}, "list"),
		},
		{
			name:     "Test with nil map",
			input:    nil,
			expected: nil, //out is interface, no []string, so zero is nil
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = []string{}
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

func TestToList_List(t *testing.T) {
	type account struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var tests = []structTest{
		{
			name:     "Test slice to slice",
			input:    []interface{}{"1", "2"},
			expected: []string{"1", "2"},
		},
		{
			name:     "Test array to slice",
			input:    [2]interface{}{"1", "2"},
			expected: []string{"1", "2"},
		},
		{
			name:     "Test slice to array",
			input:    []interface{}{1, 2},
			output:   [2]string{},
			expected: [2]string{"1", "2"},
		},
		{
			name:     "Test with list of structs",
			input:    []interface{}{&account{Id: 1, Name: "abc"}, &account{Id: 2, Name: "def"}},
			output:   [2]string{},
			expected: [2]string([2]string{"{\"id\":1,\"name\":\"abc\"}", "{\"id\":2,\"name\":\"def\"}"}),
			err:      nil,
		},
		{
			name:     "Test with list of structs to interface array",
			input:    []interface{}{&account{Id: 1, Name: "abc"}, &account{Id: 2, Name: "def"}},
			output:   [2]interface{}{},
			expected: [2]interface{}{&account{Id: 1, Name: "abc"}, &account{Id: 2, Name: "def"}},
			err:      nil,
		},
		{
			name:     "Test with nil list",
			input:    nil,
			expected: nil, //out is interface, no []string, so zero is nil
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = []string{}
			}
			if tt.op == nil {
				tt.op = NewOptions()
			}
			var err error
			err = ToAny(tt.input, &result, *tt.op)
			if err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestToList_String(t *testing.T) {
	type account struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var tests = []structTest{
		{
			name:     "Test string to slice",
			input:    `["1", "2"]`,
			output:   []string{},
			expected: []string{"1", "2"},
		},
		{
			name:     "Test nil to slice",
			input:    ``,
			output:   []string{},
			expected: []string(nil),
			err:      nil,
		},
		{
			name:     "Test string map to slice",
			input:    `{"id": 1}`,
			output:   []string{},
			expected: []string{"1"}, //this is equivalent to map to list
		},
		{
			name:     "Test string map to slice",
			input:    `[{"id": 1, "name": "abc"}, {"id": 2, "name": "def"}]`,
			output:   []*account{},
			expected: []*account{{Id: 1, Name: "abc"}, {Id: 2, Name: "def"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = tt.output
			if result == nil {
				result = []string{}
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
