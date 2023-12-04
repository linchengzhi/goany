package goany

import (
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

// supported basic type
// int, int8, int16, int32, int64
// uint, uint8, uint16, uint32, uint64
// float32, float64
// bool
// string

// ToInt64 convert an interface to an int64 type.
func ToInt64(v interface{}) int64 {
	out, _ := toInt64E(v)
	return out
}

func ToInt32(v interface{}) int32 {
	out, _ := toInt64E(v)
	return int32(out)
}

func ToInt16(v interface{}) int16 {
	out, _ := toInt64E(v)
	return int16(out)
}

func ToInt8(v interface{}) int8 {
	out, _ := toInt64E(v)
	return int8(out)
}

func ToInt(v interface{}) int {
	out, _ := toInt64E(v)
	return int(out)
}

// ToInt64E is a function that attempts to convert an arbitrary type to an int64.
func ToInt64E(v interface{}, op ...Options) (int64, error) {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	return toInt64E(v, op[0])
}

func ToUint64(v interface{}) uint64 {
	out, _ := toUint64E(v)
	return out
}

func ToUint32(v interface{}) uint32 {
	out, _ := toUint64E(v)
	return uint32(out)
}

func ToUint16(v interface{}) uint16 {
	out, _ := toUint64E(v)
	return uint16(out)
}

func ToUint8(v interface{}) uint8 {
	out, _ := toUint64E(v)
	return uint8(out)
}

func ToUint(v interface{}) uint {
	out, _ := toUint64E(v)
	return uint(out)
}

func ToUint64E(v interface{}, op ...Options) (uint64, error) {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	return toUint64E(v, op[0])
}

// ToFloat32 convert an interface to a float32 type
func ToFloat32(v interface{}) float32 {
	out, _ := toFloat64E(v)
	return float32(out)
}

// ToFloat64 convert an interface to a float64 type
func ToFloat64(v interface{}) float64 {
	out, _ := toFloat64E(v)
	return out
}

func ToFloat64E(v interface{}, op ...Options) (float64, error) {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	return toFloat64E(v, op[0])
}

// ToString convert an interface to a string type
func ToString(v interface{}) string {
	out, _ := toStringE(v)
	return out
}

func ToStringE(v interface{}, op ...Options) (string, error) {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	return toStringE(v, op[0])
}

func ToBool(v interface{}) bool {
	out, _ := toBoolE(v)
	return out
}

// ToBoolE convert an interface to a bool type
func ToBoolE(v interface{}, op ...Options) (bool, error) {
	if len(op) == 0 {
		op = append(op, *NewOptions())
	}
	return toBoolE(v, op[0])
}

// decodeBasic is a function that attempts to convert an arbitrary type to a basic type.
func (cli *anyClient) decodeBasic(in interface{}, outVal reflect.Value) error {
	switch outVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result, err := toInt64E(in, *cli.options)
		if err != nil {
			return err
		}
		reflect.NewAt(outVal.Type(), unsafe.Pointer(outVal.UnsafeAddr())).Elem().SetInt(result)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result, err := toUint64E(in)
		if err != nil {
			return err
		}
		reflect.NewAt(outVal.Type(), unsafe.Pointer(outVal.UnsafeAddr())).Elem().SetUint(result)
	case reflect.Float32, reflect.Float64:
		result, err := toFloat64E(in)
		if err != nil {
			return err
		}
		reflect.NewAt(outVal.Type(), unsafe.Pointer(outVal.UnsafeAddr())).Elem().SetFloat(result)
	case reflect.Bool:
		result, err := toBoolE(in)
		if err != nil {
			return err
		}
		reflect.NewAt(outVal.Type(), unsafe.Pointer(outVal.UnsafeAddr())).Elem().SetBool(result)
	case reflect.String:
		result, err := toStringE(in, *cli.options)
		if err != nil {
			return err
		}
		reflect.NewAt(outVal.Type(), unsafe.Pointer(outVal.UnsafeAddr())).Elem().SetString(result)
	case reflect.Ptr:
		return cli.decodeBasic(in, outVal.Elem())
	default:
		return errors.Errorf(ErrUnableConvertBasic, in)
	}
	return nil
}

// toInt64E is a function that attempts to convert an arbitrary type to an int64.
// It accepts an interface{} as input, which can be a value of any type, and optional Options.
// It returns an int64 and an error. If the conversion is successful, the error is nil.
// If the input is nil, it returns 0 and nil.
//
// The function handles the following types:
// - String: It uses strconv.ParseInt to convert the string to an int64.
// - Int, Int8, Int16, Int32, Int64: It directly returns the int value.
// - Uint, Uint8, Uint16, Uint32, Uint64: It converts the uint value to an int64 and returns it.
// - Float32, Float64: It converts the float value to an int64 and returns it.
// - Bool: It returns 1 if the bool is true, and 0 if it's false.
// - Struct: If the struct is of type time.Time, it returns the Unix timestamp. Otherwise, it returns an error.
// - Slice: If the slice is of type []uint8å (a byte slice), it converts the byte slice to a string and then uses strconv.ParseInt to convert the string to an int64. Otherwise, it returns an error.
//
// If the input is of any other type, the function returns an error.
func toInt64E(v interface{}, op ...Options) (int64, error) {
	v = Indirect(v)
	if CheckInIsNil(v) {
		return 0, nil
	}

	switch outVal := reflect.ValueOf(v); outVal.Kind() {
	case reflect.String:
		return strconv.ParseInt(outVal.String(), 10, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return outVal.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(outVal.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(outVal.Float()), nil
	case reflect.Bool:
		if outVal.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v.(time.Time).Unix(), nil
		default:
			return 0, errors.Errorf(ErrUnableConvertInt64, v)
		}
	case reflect.Slice:
		switch v.(type) {
		case []uint8:
			return strconv.ParseInt(string(v.([]uint8)), 10, 64)
		default:
			return 0, errors.Errorf(ErrUnableConvertInt64, v)
		}
	default:
		return 0, errors.Errorf(ErrUnableConvertInt64, v)
	}
}

// toUint64E converts an interface to an uint64 type.
func toUint64E(v interface{}, op ...Options) (uint64, error) {
	v = Indirect(v)
	if CheckInIsNil(v) {
		return 0, nil
	}
	switch outVal := reflect.ValueOf(v); outVal.Kind() {
	case reflect.String:
		return strconv.ParseUint(outVal.String(), 10, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(outVal.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return outVal.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(outVal.Float()), nil
	case reflect.Bool:
		if outVal.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Slice:
		switch v.(type) {
		case []uint8:
			return strconv.ParseUint(string(v.([]uint8)), 10, 64)
		default:
			return 0, errors.Errorf(ErrUnableConvertUint64, v)
		}
	default:
		return 0, errors.Errorf(ErrUnableConvertUint64, v)
	}
}

// toFloat64E converts an interface to a float64 type.
func toFloat64E(v interface{}, op ...Options) (float64, error) {
	v = Indirect(v)
	if CheckInIsNil(v) {
		return 0, nil
	}
	switch outVal := reflect.ValueOf(v); outVal.Kind() {
	case reflect.String:
		d, err := strconv.ParseFloat(outVal.String(), 64)
		return d, err
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(outVal.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(outVal.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return outVal.Float(), nil
	case reflect.Bool:
		if outVal.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Slice:
		switch v.(type) {
		case []uint8:
			return strconv.ParseFloat(string(v.([]uint8)), 64)
		default:
			return 0, errors.Errorf(ErrUnableConvertFloat64, v)
		}
	default:
		return 0, errors.Errorf(ErrUnableConvertFloat64, v)
	}
}

// toStringE converts an interface to a string type.
func toStringE(v interface{}, op ...Options) (string, error) {
	v = Indirect(v)
	if CheckInIsNil(v) {
		return "", nil
	}
	switch outVal := reflect.ValueOf(v); outVal.Kind() {
	case reflect.String:
		if _, ok := v.(string); ok {
			return v.(string), nil
		}
		return outVal.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(outVal.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(outVal.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(outVal.Float(), 'g', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(v.(bool)), nil
	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v.(time.Time).Format(op[0].timeFormat), nil
		default:
			b, err := json.Marshal(v)
			return string(b), err
		}
	case reflect.Map:
		b, err := json.Marshal(v)
		return string(b), err
	case reflect.Slice, reflect.Array:
		switch v.(type) {
		case []uint8:
			return string(v.([]uint8)), nil
		default:
			b, err := json.Marshal(v)
			return string(b), err
		}
	default:
		return "", errors.Errorf(ErrUnableConvertString, v)
	}
}

// toBoolE converts an interface to a bool type.
func toBoolE(v interface{}, op ...Options) (bool, error) {
	v = Indirect(v)
	if CheckInIsNil(v) {
		return false, nil
	}
	switch outVal := reflect.ValueOf(v); outVal.Kind() {
	case reflect.String:
		return strconv.ParseBool(outVal.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return outVal.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return outVal.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return outVal.Float() != 0, nil
	case reflect.Bool:
		return outVal.Bool(), nil
	case reflect.Slice:
		switch v.(type) {
		case []uint8:
			return strconv.ParseBool(string(v.([]uint8)))
		default:
			return false, errors.Errorf(ErrUnableConvertBool, v)
		}
	default:
		return false, errors.Errorf(ErrUnableConvertBool, v)
	}
}
