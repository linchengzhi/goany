package goany

import (
	"github.com/pkg/errors"
	"reflect"
)

// decodeMap decodes an input value into a map output value. The input can be a map, struct, list (array or slice),
// string (in JSON format), or pointer. The appropriate decoding method is chosen based on the input kind.
func (cli *anyClient) decodeMap(in interface{}, outVal reflect.Value) error {
	inVal := reflect.Indirect(reflect.ValueOf(in))

	switch inVal.Kind() {
	case reflect.Map:
		return cli.mapToMap(in, outVal)
	case reflect.Struct:
		return cli.structToMap(in, outVal)
	case reflect.Array, reflect.Slice:
		return cli.listToMap(in, outVal)
	case reflect.String:
		return cli.stringToAny(in, outVal)
	case reflect.Ptr:
		return cli.decodeAny(inVal.Elem().Interface(), outVal)
	default:
		return errors.Errorf(ErrInToOut, in, "map")
	}
}

// mapToMap converts a map input into a map output value. It decodes each key-value pair from the input map
// and sets the corresponding key-value pair in the output map. The output map's key and value types are determined
// dynamically at runtime.
func (cli *anyClient) mapToMap(in interface{}, outVal reflect.Value) error {
	basicOutVal := reflect.MakeMap(outVal.Type())
	basicOutKey := basicOutVal.Type().Key()
	basicOutElem := basicOutVal.Type().Elem()

	inVal := reflect.ValueOf(in)

	for _, k := range inVal.MapKeys() {
		currentKey := reflect.Indirect(reflect.New(basicOutKey))
		if err := cli.decodeAny(k.Interface(), currentKey); err != nil {
			return err
		}
		inFieldVal := inVal.MapIndex(k).Interface()
		currentValue := reflect.Indirect(reflect.New(basicOutElem))

		if err := cli.decodeAny(inFieldVal, currentValue); err != nil {
			return err
		}
		basicOutVal.SetMapIndex(currentKey, currentValue)
	}
	outVal.Set(basicOutVal)
	return nil
}

// structToMap converts a struct input into a map output value. Each exported field of the struct (or unexported if
// the `exportedLower` option is set) is decoded and added to the map using the field's name as the key. Fields tagged
// with `TagIgnore` are skipped. If `exportedLower` is not set and the field is unexported, it is ignored unless it
// is accessible via a pointer.
func (cli *anyClient) structToMap(in interface{}, outVal reflect.Value) error {
	basicOutVal := reflect.MakeMap(outVal.Type())
	basicOutKey := basicOutVal.Type().Key()
	basicOutElem := basicOutVal.Type().Elem()

	inType, inValue := ReflectTypeValue(in)

	for i := 0; i < inType.NumField(); i++ {
		inField := new(fieldInfo)
		inField.fieldStruct = inType.Field(i)

		currentKey := reflect.Indirect(reflect.New(basicOutKey))
		currentValue := reflect.Indirect(reflect.New(basicOutElem))
		inField.fieldName = GetFieldNameByTag(inField.fieldStruct, cli.options.tagName)
		if !inField.canUse(*cli.options) {
			continue
		}

		if err := cli.decodeAny(inField.fieldName, currentKey); err != nil {
			return err
		}

		inFieldVal := inValue.Field(i)
		if inField.fieldStruct.PkgPath != "" { //if inField is unexported, it must be accessible via a pointer
			if !inFieldVal.CanAddr() {
				return ErrInNotPtr
			}
			inFieldVal = getUnexportedField(inFieldVal)
		}

		if err := cli.decodeAny(inFieldVal.Interface(), currentValue); err != nil {
			return err
		}
		basicOutVal.SetMapIndex(currentKey, currentValue)
	}
	outVal.Set(basicOutVal)
	return nil
}

// listToMap converts a list input (array or slice) into a map output value. If the `mapKeyField` option is set and the
// elements are structs, the specified field within each struct is used as the map key. Otherwise, the list index is used
// as the map key. Each element of the list is decoded and added to the map.
func (cli *anyClient) listToMap(in interface{}, outVal reflect.Value) error {
	basicOutVal := reflect.MakeMap(outVal.Type())
	basicOutKey := basicOutVal.Type().Key()
	basicOutElem := basicOutVal.Type().Elem()

	_, inValue := ReflectTypeValue(in)

	for i := 0; i < inValue.Len(); i++ {
		currentKey := reflect.Indirect(reflect.New(basicOutKey))
		currentValue := reflect.Indirect(reflect.New(basicOutElem))

		inFiledVal := inValue.Index(i)
		inFiledValElem := reflect.Indirect(inFiledVal)

		// if mapKeyField is set, use the specified field as the map key
		if cli.options.mapKeyField != "" && inFiledValElem.Kind() == reflect.Struct {
			var keyFieldVal reflect.Value
			for j := 0; j < inFiledValElem.NumField(); j++ {
				fileName := GetFieldNameByTag(inFiledValElem.Type().Field(j), cli.options.tagName)
				if fileName == cli.options.mapKeyField {
					keyFieldVal = inFiledValElem.Field(j)
					break
				}
			}
			if keyFieldVal.Kind() == reflect.Invalid {
				return errors.Errorf(ErrFieldNoFound, cli.options.mapKeyField)
			}
			if err := cli.decodeAny(keyFieldVal.Interface(), currentKey); err != nil {
				return err
			}
		} else {
			if err := cli.decodeAny(i, currentKey); err != nil {
				return err
			}
		}
		if err := cli.decodeAny(inFiledVal.Interface(), currentValue); err != nil {
			return err
		}
		basicOutVal.SetMapIndex(currentKey, currentValue)
	}
	outVal.Set(basicOutVal)
	return nil
}
