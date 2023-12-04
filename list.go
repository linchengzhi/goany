package goany

import (
	"github.com/pkg/errors"
	"reflect"
)

// decodeList decodes an input value into a list (slice or array) output value.
// It supports decoding from various input types such as maps, arrays, slices, strings (in JSON format), and pointers.
func (cli *anyClient) decodeList(in interface{}, outVal reflect.Value) error {
	inVal := reflect.Indirect(reflect.ValueOf(in))
	switch inVal.Kind() {
	case reflect.Map:
		// Decode a map into a list by iterating over its elements.
		return cli.mapToList(in, outVal)
	case reflect.Array, reflect.Slice:
		// Decode an input list into an output list element by element.
		return cli.listToList(in, outVal)
	case reflect.String:
		// If the input is a string, attempt to decode it as JSON into a list.
		return cli.stringToAny(in, outVal)
	case reflect.Ptr:
		return cli.decodeAny(inVal.Elem().Interface(), outVal)
	default:
		return errors.Errorf(ErrInToOut, in, "list")
	}
}

// mapToList converts a map input into a slice output value. If the mapKeyToList option is true,
// the map's keys are used to populate the slice, otherwise the map's values are used.
func (cli *anyClient) mapToList(in interface{}, outVal reflect.Value) error {
	_, inVal := ReflectTypeValue(in)
	var basicOutVal reflect.Value
	if outVal.Kind() == reflect.Array {
		// If the output is an array, create a new array of the appropriate type and size.
		arrType := reflect.ArrayOf(outVal.Len(), outVal.Type().Elem())
		basicOutVal = reflect.New(arrType).Elem()
	} else {
		// If the output is a slice, create a new slice of the appropriate type and size.
		basicOutVal = reflect.MakeSlice(outVal.Type(), inVal.Len(), inVal.Len())
	}
	iter := inVal.MapRange()
	i := 0
	for iter.Next() {
		v := iter.Value().Interface()
		if cli.options.mapKeyToList {
			v = iter.Key().Interface()
		}
		if err := cli.decodeAny(v, basicOutVal.Index(i)); err != nil {
			return err
		}
		i++
	}
	outVal.Set(basicOutVal)
	return nil
}

// listToList converts an input list (array or slice) into an output list.
// It decodes each element from the input list and sets it in the output list.
func (cli *anyClient) listToList(in interface{}, outVal reflect.Value) error {
	_, inVal := ReflectTypeValue(in)
	var basicOutVal reflect.Value
	if outVal.Kind() == reflect.Array {
		// If the output is an array, create a new array of the appropriate type and size.
		arrType := reflect.ArrayOf(outVal.Len(), outVal.Type().Elem())
		basicOutVal = reflect.New(arrType).Elem()
	} else {
		// If the output is a slice, create a new slice of the appropriate type and size.
		basicOutVal = reflect.MakeSlice(outVal.Type(), inVal.Len(), inVal.Len())
	}
	for i := 0; i < inVal.Len(); i++ {
		v := inVal.Index(i).Interface()
		if err := cli.decodeAny(v, basicOutVal.Index(i)); err != nil {
			return err
		}
	}
	outVal.Set(basicOutVal)
	return nil
}
