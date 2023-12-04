package goany

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToTime(t *testing.T) {
	locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
	input := "2020-10-01T21:06:11"
	expected := time.Date(2020, 10, 1, 21, 6, 11, 0, locationShanghai)
	op := NewOptions().SetLocation(locationShanghai)
	actual := ToTime(input, *op)
	assert.Equal(t, expected.Format("2006-01-02T15:04:05Z07:00"), actual.Format("2006-01-02T15:04:05Z07:00"))
}

func TestToTimeE(t *testing.T) {
	locationUTC, _ := time.LoadLocation("UTC")
	locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
	locationBarnaul, _ := time.LoadLocation("Asia/Barnaul") //+07:00
	fmt.Println(locationBarnaul)

	tests := []struct {
		name     string
		input    interface{}
		op       Options
		expected time.Time
		err      error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: time.Time{},
			err:      nil,
		},
		{
			name:     "time input",
			input:    time.Date(2023, 11, 4, 21, 6, 11, 0, locationUTC),
			expected: time.Date(2023, 11, 4, 21, 6, 11, 0, locationUTC),
			err:      nil,
		},
		{
			name:     "string input",
			input:    "2020-10-01T21:06:11Z",
			expected: time.Date(2020, 10, 1, 21, 6, 11, 0, locationUTC),
			err:      nil,
		},
		{
			name:     "string input, loc is +7",
			input:    "2020-01-01T21:30:03+07:00",
			expected: time.Date(2020, 1, 1, 21, 30, 03, 0, locationBarnaul),
			err:      nil,
		},
		{
			name:     "invalid string input",
			input:    "invalid",
			expected: time.Time{},
			err:      errors.Errorf(ErrUnableConvertTime, "invalid"),
		},
		{
			name:     "integer input",
			input:    int64(1636058771), //2021-11-05 04:46:11 shanghai
			op:       *NewOptions().SetLocation(locationShanghai),
			expected: time.Date(2021, 11, 5, 4, 46, 11, 0, locationShanghai),
			err:      nil,
		},
		{
			name:     "int ptr input, loc is shanghai",
			input:    func() *int64 { i := int64(1636058771); return &i }(), //2021-11-05 04:46:11
			expected: time.Date(2021, 11, 4, 20, 46, 11, 0, locationUTC),
			err:      nil,
		},
		{
			name:     "invalid type input",
			input:    []int{1, 2, 3},
			expected: time.Time{},
			err:      errors.Errorf(ErrUnableConvertTime, []int{1, 2, 3}),
		},
		{
			name:     "time input with options",
			input:    time.Date(2023, 11, 4, 21, 6, 11, 0, locationShanghai),
			op:       *NewOptions().SetLocation(locationShanghai),
			expected: time.Date(2023, 11, 4, 21, 6, 11, 0, locationShanghai),
		},
		{
			name:     "null int ptr input",
			input:    func() *int64 { return nil }(),
			expected: time.Time{},
			err:      nil,
		},
		{
			name:     "struct input",
			input:    struct{}{},
			expected: time.Time{},
			err:      errors.Errorf(ErrUnableConvertTime, struct{}{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ToTimeE(tt.input, tt.op)
			assert.Equal(t, tt.expected.Format("2006-01-02T15:04:05Z07:00"),
				actual.Format("2006-01-02T15:04:05Z07:00"))
			if tt.err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			}
		})
	}
}

func TestToTime_ToAny(t *testing.T) {
	locationUTC, _ := time.LoadLocation("UTC")

	tests := []struct {
		name     string
		input    interface{}
		options  []Options
		expected time.Time
		err      error
	}{
		{
			name:     "time input",
			input:    time.Date(2023, 11, 4, 21, 6, 11, 0, locationUTC),
			expected: time.Date(2023, 11, 4, 21, 6, 11, 0, locationUTC),
			err:      nil,
		},
		{
			name:     "string input",
			input:    "2020-10-01T21:06:11Z",
			expected: time.Date(2020, 10, 1, 21, 6, 11, 0, locationUTC),
			err:      nil,
		},
		{
			name:     "integer input",
			input:    int64(1636058771), //2021-11-05 04:46:11
			expected: time.Date(2021, 11, 4, 20, 46, 11, 0, locationUTC),
			err:      nil,
		},
		{
			name:     "err input",
			input:    []int{1, 2, 3},
			expected: time.Date(2021, 11, 4, 20, 46, 11, 0, locationUTC),
			err:      errors.Errorf(ErrUnableConvertTime, []int{1, 2, 3}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out time.Time
			err := ToAny(tt.input, &out, tt.options...)
			if err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, out)
			}
			fmt.Println(out.Format("2006-01-02 15:04:05"))
		})
	}
}

func TestTimeToString(t *testing.T) {
	locationShanghai, _ := time.LoadLocation("Asia/Shanghai")
	in := time.Date(2020, 10, 1, 21, 6, 11, 0, locationShanghai)
	expected := "2020-10-01T21:06:11+08:00"
	op := NewOptions().SetLocation(locationShanghai).SetTimeFormat(time.RFC3339)
	out, err := toStringE(in, *op)
	fmt.Println(out, err)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
