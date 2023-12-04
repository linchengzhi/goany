package goany

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"time"
)

// IAny defines an interface for converting various types to the 'any' type.
type IAny interface {
	// decodeAny is used to decode an input interface to the desired output value.
	decodeAny(in interface{}, outVal reflect.Value) error

	// The following methods are specialized decoders for different kinds of values.
	decodeBasic(in interface{}, outVal reflect.Value) error

	decodeInterface(in interface{}, outVal reflect.Value) error
	decodePtr(in interface{}, outVal reflect.Value) error

	decodeMap(in interface{}, outVal reflect.Value) error
	decodeStruct(in interface{}, outVal reflect.Value) error
	decodeList(in interface{}, outVal reflect.Value) error

	decodeTime(in interface{}, outVal reflect.Value) error
}

// ToAny converts an interface to an any type. It uses reflection to dynamically
// decode the input value into the output container based on the output's type.
func ToAny(in interface{}, out interface{}, options ...Options) error {

	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr { // if out is not ptr, return error
		return errors.Errorf(ErrUnSupportType, out)
	}

	outVal = outVal.Elem()
	if !outVal.CanAddr() { // if out can not addr, return error
		return errors.Errorf(ErrUnSupportType, out)
	}

	cli := newAnyClient(options...) // new a client, init options

	err := cli.decodeAny(in, outVal)
	if err != nil && err != ErrDecodeStop {
		return err
	}
	return nil
}

// decodeAny attempts to decode the input value into the provided output value.
// It determines how to decode the input based on the type of the output value.
func (cli *anyClient) decodeAny(in interface{}, outVal reflect.Value) error {

	// If the input value is nil or a nil pointer, set the output to its zero value.
	if CheckInIsNil(Indirect(in)) {
		outVal.Set(reflect.Zero(outVal.Type()))
		return nil
	}

	// If there are decoding hooks defined, process them.
	if len(cli.options.hooks) > 0 {
		// Execute the hook, and if it returns an error, stop the process.
		for _, hook := range cli.options.hooks {
			status, err := hook(in, outVal)
			if err != nil {
				return err
			}
			if status == DecodeSkip { // if hook return skip, skip decode
				return nil
			}
			if status == DecodeStop { // if hook return stop, stop decode
				return ErrDecodeStop
			}
		}
	}

	// Based on the kind of the output value, call the appropriate decoding function.
	outKind := outVal.Kind()
	var err error
	switch {
	case isBasicType(outKind):
		err = cli.decodeBasic(in, outVal)
	case outKind == reflect.Interface:
		err = cli.decodeInterface(in, outVal)
	case outKind == reflect.Map:
		err = cli.decodeMap(in, outVal)
	case outKind == reflect.Struct:
		// Special case for time.Time type which requires specific handling.
		switch outVal.Interface().(type) {
		case time.Time:
			err = cli.decodeTime(in, outVal)
		default:
			err = cli.decodeStruct(in, outVal)
		}
	case outKind == reflect.Slice || outKind == reflect.Array:
		err = cli.decodeList(in, outVal)
	case outKind == reflect.Ptr:
		err = cli.decodePtr(in, outVal)
	default:
		err = errors.Errorf(ErrUnSupportType, outKind.String())
	}
	if err != nil {
		return err
	}
	return nil
}

// decodeInterface handles the decoding of an interface value into the provided output value.
// The function works as follows:
//  1. If the output value (outVal) is valid and not nil, it creates a new instance of the type
//     that outVal represents, decodes the input into it, and sets the outVal with the new instance.
//  2. If the options are set to export detailed information and the input is a struct,
//     it converts the input struct into a map with string keys and interface{} values,
//     allowing for more detailed introspection of struct fields.
//  3. If the output value is assignable from the input value's type, it directly assigns the input to the output.
//  4. If the output value is not assignable, it returns an error indicating that the
//     input value cannot be assigned to the output interface.
func (cli *anyClient) decodeInterface(in interface{}, outVal reflect.Value) error {
	// Check if outVal is valid and has a non-nil underlying value.
	if outVal.IsValid() && outVal.Elem().IsValid() {
		outValElem := outVal.Elem()
		currentOutVal := reflect.New(outValElem.Type()).Elem()
		if err := cli.decodeAny(in, currentOutVal); err != nil {
			return err
		}
		outVal.Set(currentOutVal)
		return nil
	}

	inVal := reflect.ValueOf(in)
	inValElem := reflect.Indirect(inVal)

	// If exporting detailed struct information is enabled and the input is a struct,
	// convert it into a map for detailed introspection.
	if cli.options.structToMapDetail && inValElem.Kind() == reflect.Struct {
		currentMap := make(map[string]interface{})
		currentValue := reflect.ValueOf(&currentMap).Elem()
		if err := cli.decodeAny(in, currentValue); err != nil {
			return err
		}
		outVal.Set(currentValue)
	} else {
		// Check if the input value's type is assignable to the output value's type.
		if !inVal.Type().AssignableTo(outVal.Type()) {
			return errors.Errorf(ErrInToOut, in, "interface")
		}
		outVal.Set(inVal)
	}
	return nil
}

// decodePtr handles the decoding of pointer types. It takes an input value and a reflect.Value
// that represents a pointer. The method first checks if the output value (outVal) is nil. If it is,
// it creates a new instance of the type that the pointer points to, decodes the input value into
// that new instance, and then sets the original pointer to point to this new instance.
// If the output value is not nil, it decodes the input value into the value that the pointer
// is already pointing to.
func (cli *anyClient) decodePtr(in interface{}, outVal reflect.Value) error {
	if outVal.IsNil() {
		var newVal = reflect.New(outVal.Type().Elem())
		if err := cli.decodeAny(in, reflect.Indirect(newVal)); err != nil {
			return err
		}
		outVal.Set(newVal)
	} else {
		if err := cli.decodeAny(in, reflect.Indirect(outVal)); err != nil {
			return err
		}
	}
	return nil
}

func (cli *anyClient) stringToAny(in interface{}, outVal reflect.Value) error {
	inBytes := []byte(reflect.ValueOf(in).String())
	// check if the in is valid json
	if !json.Valid(inBytes) {
		return errors.Errorf(ErrNotJson, in)
	}
	var inData interface{}
	dec := json.NewDecoder(bytes.NewBuffer(inBytes))
	dec.UseNumber()
	err := dec.Decode(&inData)
	if err != nil {
		return errors.Errorf(ErrNotJson, in)
	}

	inDataType, _ := ReflectTypeValue(inData)
	switch inDataType.Kind() {
	case reflect.Map, reflect.Slice:
		return cli.decodeAny(inData, outVal)
	default:
		return errors.Errorf(ErrNotJson, in) // reject other types
	}
}

// isBasicType returns true if the type is a basic type.
func isBasicType(v reflect.Kind) bool {
	switch v {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.String:
	default:
		return false
	}
	return true
}

// check in is nil
func CheckInIsNil(in interface{}) bool {
	var inVal reflect.Value
	if in != nil {
		inVal = reflect.ValueOf(in)
		if inVal.Kind() == reflect.Ptr && inVal.IsNil() {
			in = nil
		}
	}
	return in == nil
}

// The Indirect function dereferences an interface{} value
// until it reaches a non-pointer base value or a nil value.
func Indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// ReflectTypeValue returns the reflect.Type and reflect.Value of source.
func ReflectTypeValue(source interface{}) (reflect.Type, reflect.Value) {
	if source == nil {
		return nil, reflect.ValueOf(nil)
	}
	sourceRt := reflect.TypeOf(source)
	sourceRv := reflect.ValueOf(source)
	for sourceRt.Kind() == reflect.Ptr && !sourceRv.IsNil() {
		sourceRt = sourceRt.Elem()
		sourceRv = sourceRv.Elem()
	}
	return sourceRt, sourceRv
}

// GetFieldNameByTag returns the field name by tag.
// if tag is gorm, get column name
// if tag name is nil, get field name
func GetFieldNameByTag(field reflect.StructField, tag string) string {
	fieldName := field.Tag.Get(tag)
	if tag == "gorm" && strings.Contains(fieldName, "column:") {
		fieldName = strings.Split(fieldName, "column:")[1]
	}
	if strings.Contains(fieldName, ",") {
		fieldName = strings.Split(fieldName, ",")[0]
	}
	if fieldName == "" { //no tag, get field name
		fieldName = field.Name
	}
	return fieldName
}
