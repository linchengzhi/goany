package goany

import (
	"github.com/pkg/errors"
	"reflect"
	"unsafe"
)

// decodeStruct decodes an input value into a struct output value. It handles different input types
// like maps, structs, strings (in JSON format), and pointers. Depending on the input kind, it delegates
// to the corresponding method for decoding.
func (cli *anyClient) decodeStruct(in interface{}, outVal reflect.Value) error {
	inVal := reflect.Indirect(reflect.ValueOf(in))

	switch inVal.Kind() {
	case reflect.Map:
		return cli.mapToStruct(in, outVal)
	case reflect.Struct:
		return cli.structToStruct(in, outVal)
	case reflect.String:
		// If the input is a string, attempt to decode it as JSON.
		return cli.stringToAny(in, outVal)
	case reflect.Ptr:
		return cli.decodeAny(inVal.Elem().Interface(), outVal)
	default:
		// If none of the above types match, return an error indicating unsupported type.
		return errors.Errorf(ErrInToOut, in, "struct")
	}
}

// mapToStruct decodes a map input into a struct output value. It iterates over each field
// in the input map and attempts to match and set the corresponding field in the output struct.
func (cli *anyClient) mapToStruct(in interface{}, outVal reflect.Value) error {
	// Create a new instance of the output value's type.
	basicOutVal := reflect.New(outVal.Type())
	basicOutValElem := basicOutVal.Elem()

	_, inVal := ReflectTypeValue(in)

	// Extract field information from the input map and output struct.
	inFieldInfos := deepMapInFields(inVal)
	outFieldInfos, outAnonymous := deepOutFields(basicOutVal, "", *cli.options)
	for _, inFieldInfo := range inFieldInfos {

		// Skip the field if the input map key does not have a corresponding field in the output struct
		matchOuts := matchOutField(inFieldInfo.fieldName, outFieldInfos, cli.options.assignKey)
		for _, outFieldInfo := range matchOuts {
			if outFieldInfo == nil {
				continue
			}

			if err := cli.decodeAny(inFieldInfo.fieldVal.Interface(), outFieldInfo.fieldVal); err != nil {
				return err
			}
			// Remove the field from the map of output fields to avoid multiple assignments.
			delete(outFieldInfos, outFieldInfo.fieldName)

			for _, anon := range outAnonymous[outFieldInfo.fieldName] {
				delete(outFieldInfos, anon)
			}
		}
	}
	//if outfield not match, set it to nil
	for _, v := range outFieldInfos {
		if v.fieldStruct.Anonymous {
			continue
		}
		v.fieldVal.Set(reflect.Zero(v.fieldStruct.Type))
	}
	outVal.Set(basicOutValElem)
	return nil
}

// structToStruct decodes a struct input into a struct output value. It handles matching
// fields between the input and output structs and decoding each field.
func (cli *anyClient) structToStruct(in interface{}, outVal reflect.Value) error {
	// Create a new instance of the output value's type.
	basicOutVal := reflect.New(outVal.Type()).Elem()

	inFieldInfos, err := deepInFields(in, "", *cli.options)
	if err != nil {
		return err
	}
	outFieldInfos, outAnonymous := deepOutFields(basicOutVal, "", *cli.options)

	var currentAnonymous = ""
	for _, inFieldInfo := range inFieldInfos {
		if inFieldInfo.belongAnonymous != "" && inFieldInfo.belongAnonymous == currentAnonymous {
			// Skip fields that belong to an already handled anonymous struct.
			continue
		}

		// If no matching field is found, skip to the next field.
		matchOuts := matchOutField(inFieldInfo.fieldName, outFieldInfos, cli.options.assignKey)
		for _, outFieldInfo := range matchOuts {
			if err := cli.decodeAny(inFieldInfo.fieldVal.Interface(), outFieldInfo.fieldVal); err != nil {
				return err
			}

			// Remove the field from the map of output fields to avoid multiple assignments.
			delete(outFieldInfos, outFieldInfo.fieldName)

			// Also, remove any Anonymous associated with the field.
			for _, anon := range outAnonymous[outFieldInfo.fieldName] {
				delete(outFieldInfos, anon)
			}
		}

		// If the current field is an anonymous struct, set it as currentAnonymous.
		if inFieldInfo.isAnonymous {
			currentAnonymous = inFieldInfo.fieldName
		}
	}
	//if outfield not match, set it to nil
	for _, v := range outFieldInfos {
		if v.fieldStruct.Anonymous { //anonymous field can not be set to nil
			continue
		}
		v.fieldVal.Set(reflect.Zero(v.fieldStruct.Type))
	}

	outVal.Set(basicOutVal)
	return nil
}

// matchOutField attempts to find a field in the output struct that matches the input field name.
// It takes into consideration any custom assignKey mappings that may be used to match fields
// with different names between the input and output.
func matchOutField(inKey string, outKeys map[string]fieldInfo, assignKey map[string]string) []*fieldInfo {
	// Check if there is a direct match for the input key in the output keys.
	var list = make([]*fieldInfo, 0)
	if out, ok := outKeys[inKey]; ok {
		list = append(list, &out)
	}

	// Check if there is an assignKey mapping for the input key and if it matches an output key.
	if key, ok := assignKey[inKey]; ok {
		if out, ok := outKeys[key]; ok {
			list = append(list, &out)
		}
	}
	return list
}

// get unexported field value
func getUnexportedField(field reflect.Value) reflect.Value {
	if field.Kind() == reflect.Ptr {
		return reflect.NewAt(field.Type().Elem(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	} else {
		return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	}
}
