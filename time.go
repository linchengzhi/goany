package goany

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

var (
	timeType = []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime, //2006-01-02 15:04:05
		time.DateOnly,
		time.TimeOnly,
		"2006/01/02 15:04:05",
		"2006-01-02T15:04:05",
	}
)

// ToTime attempts to convert an interface value to a time.Time value
func ToTime(v interface{}, op ...Options) time.Time {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	t, _ := toTimeLocationE(v, op[0].location)
	return t
}

// ToTimeE attempts to convert an interface value to a time.Time value.
func ToTimeE(v interface{}, op ...Options) (time.Time, error) {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	return toTimeLocationE(v, op[0].location)
}

// decodeTime decodes an input value into a time.Time output value
func (cli *anyClient) decodeTime(in interface{}, outVal reflect.Value) error {
	t, err := toTimeLocationE(in, cli.options.location)
	if err != nil {
		return err
	}
	outVal.Set(reflect.ValueOf(t))
	return nil
}

// toTimeLocationE converts an interface value to time.Time considering the provided location.
// It handles struct types (specifically time.Time), strings in various date-time formats, and
// numeric types representing Unix time. The function returns the converted time.Time value and
// an error if the conversion fails. If the location is nil, UTC is used as the default.
func toTimeLocationE(in interface{}, location *time.Location) (time.Time, error) {
	if location == nil {
		location = time.UTC
	}
	in = Indirect(in)
	if CheckInIsNil(in) {
		return time.Time{}, nil
	}
	switch reflect.TypeOf(in).Kind() {
	case reflect.Struct:
		switch in.(type) {
		case time.Time:
			return in.(time.Time), nil
		default:
			return time.Time{}, errors.Errorf(ErrUnableConvertTime, in)
		}
	case reflect.String:
		return stringToTime(in.(string), location)
	case reflect.Int, reflect.Int32, reflect.Int64:
		return time.Unix(reflect.ValueOf(in).Int(), 0).In(location), nil
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return time.Unix(int64(reflect.ValueOf(in).Uint()), 0).In(location), nil
	default:
		return time.Time{}, errors.Errorf(ErrUnableConvertTime, in)
	}
}

// Compatible with time conversion in many formats
func stringToTime(v string, location *time.Location) (time.Time, error) {
	for _, value := range timeType {
		t, err := time.ParseInLocation(value, v, location)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.Errorf(ErrUnableConvertTime, v)
}
