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
	basicOutVal := reflect.New(outVal.Type()).Elem()
	_, inVal := ReflectTypeValue(in)

	// Extract field information from the input map and output struct.
	inFieldInfos := deepMapInFields(inVal)
	outFieldInfos, outAnonymous := deepOutFields(basicOutVal, *cli.options)

	for _, inFieldInfo := range inFieldInfos {

		// Skip the field if the input map key does not have a corresponding field in the output struct
		outFieldInfo := matchOutField(inFieldInfo.fieldName, outFieldInfos, cli.options.assignKey)
		if outFieldInfo == nil {
			continue
		}

		if err := cli.decodeAny(inFieldInfo.fieldVal.Interface(), outFieldInfo.fieldVal); err != nil {
			return err
		}
		// Remove the field from the map of output fields to avoid multiple assignments.
		delete(outFieldInfos, inFieldInfo.fieldName)

		for _, anon := range outAnonymous[inFieldInfo.fieldName] {
			delete(outFieldInfos, anon)
		}
	}
	outVal.Set(basicOutVal)
	return nil
}

// structToStruct decodes a struct input into a struct output value. It handles matching
// fields between the input and output structs and decoding each field.
func (cli *anyClient) structToStruct(in interface{}, outVal reflect.Value) error {
	// Create a new instance of the output value's type.
	basicOutVal := reflect.New(outVal.Type()).Elem()
	_, inVal := ReflectTypeValue(in)

	inFieldInfos := deepInFields(inVal, "", *cli.options)
	outFieldInfos, outAnonymous := deepOutFields(basicOutVal, *cli.options)

	var currentAnonymous = ""
	for _, inFieldInfo := range inFieldInfos {
		if inFieldInfo.belongAnonymous != "" && inFieldInfo.belongAnonymous == currentAnonymous {
			// Skip fields that belong to an already handled anonymous struct.
			continue
		}

		// If no matching field is found, skip to the next field.
		outFieldInfo := matchOutField(inFieldInfo.fieldName, outFieldInfos, cli.options.assignKey)
		if outFieldInfo == nil {
			continue
		}

		if err := cli.decodeAny(inFieldInfo.fieldVal.Interface(), outFieldInfo.fieldVal); err != nil {
			return err
		}

		// Remove the field from the map of output fields to avoid multiple assignments.
		delete(outFieldInfos, inFieldInfo.fieldName)

		// Also, remove any Anonymous associated with the field.
		for _, anon := range outAnonymous[inFieldInfo.fieldName] {
			delete(outFieldInfos, anon)
		}

		// If the current field is an anonymous struct, set it as currentAnonymous.
		if inFieldInfo.isAnonymous {
			currentAnonymous = inFieldInfo.fieldName
		}
	}
	outVal.Set(basicOutVal)
	return nil
}

// matchOutField attempts to find a field in the output struct that matches the input field name.
// It takes into consideration any custom assignKey mappings that may be used to match fields
// with different names between the input and output.
func matchOutField(inKey string, outKeys map[string]fieldInfo, assignKey map[string]string) *fieldInfo {
	if inKey == TagIgnore {
		return nil
	}

	// Check if there is a direct match for the input key in the output keys.
	if out, ok := outKeys[inKey]; ok {
		return &out
	}

	// Check if there is an assignKey mapping for the input key and if it matches an output key.
	if key, ok := assignKey[inKey]; ok {
		if out, ok := outKeys[key]; ok {
			return &out
		}
	}

	return nil
}

type fieldInfo struct {
	fieldName       string
	fieldStruct     reflect.StructField
	fieldVal        reflect.Value
	isAnonymous     bool
	belongAnonymous string
}

// deepMapInFields extracts a slice of fieldInfo structs representing the key-value pairs
// present in the provided map. It iterates over each entry in the map and constructs a fieldInfo
// for it unless the key corresponds to a field that should be ignored (indicated by TagIgnore).
func deepMapInFields(reflectValue reflect.Value) []fieldInfo {
	fields := make([]fieldInfo, 0)
	iter := reflectValue.MapRange()
	for iter.Next() {
		field := new(fieldInfo)
		field.fieldVal = iter.Value()
		field.fieldName = iter.Key().String()

		if field.fieldName == TagIgnore {
			continue
		}
		fields = append(fields, *field)
	}
	return fields
}

// deepInFields recursively extracts information from a struct's fields and from fields of embedded structs.
// It creates a slice of fieldInfo structs that hold the field values and names derived from tags.
// The function handles unexported fields if the exportedLower option is set, and includes a flag to indicate
// if a field is part of an embedded struct.
func deepInFields(reflectValue reflect.Value, anonymous string, op Options) []fieldInfo {
	reflectType, _ := ReflectTypeValue(reflectValue.Interface())
	fields := make([]fieldInfo, 0, reflectType.NumField())
	for i := 0; i < reflectType.NumField(); i++ {

		field := new(fieldInfo)
		fieldStruct := reflectType.Field(i)
		field.fieldVal = reflectValue.Field(i)
		field.fieldName = GetFieldNameByTag(fieldStruct, op.tagName)
		if field.fieldName == TagIgnore {
			// Skip the field if it's tagged to be ignored.
			continue
		}

		// If the field is an anonymous struct, use unsafe to get val
		if fieldStruct.PkgPath != "" {
			if !op.exportedUnExported {
				continue
			}
			field.fieldVal = getUnexportedField(field.fieldVal)
		}
		// Retrieve unexported field values if necessary.
		field.belongAnonymous = anonymous
		fields = append(fields, *field)
		if fieldStruct.Anonymous && anonymous == "" { //anonymous field, only first level
			fields = append(fields, deepInFields(field.fieldVal, field.fieldName, op)...)
		}
	}
	return fields
}

// deepOutFields creates a mapping of field names to fieldInfo structs and a mapping of
// embedded struct names to their respective fields' names for an output struct. It processes
// each field in the struct, including those in embedded structs, and handles unexported fields
// if the exportedLower option is set.
func deepOutFields(reflectVal reflect.Value, op Options) (map[string]fieldInfo, map[string][]string) {
	fields := make(map[string]fieldInfo)
	anonymous := make(map[string][]string)
	for i := 0; i < reflectVal.Type().NumField(); i++ {
		field := new(fieldInfo)
		fieldStruct := reflectVal.Type().Field(i)
		field.fieldVal = reflectVal.Field(i)
		field.fieldName = GetFieldNameByTag(fieldStruct, op.tagName)
		if field.fieldName == TagIgnore {
			continue
		}
		// If the field is an anonymous struct, use unsafe to get val
		if fieldStruct.PkgPath != "" {
			if !op.exportedUnExported {
				continue
			}
			field.fieldVal = getUnexportedField(field.fieldVal)
		}

		// Add or update the fieldInfo in the fields mapping.
		if f, ok := fields[field.fieldName]; !ok || f.belongAnonymous != "" {
			fields[field.fieldName] = *field
		}

		// Handle embedded (anonymous) struct fields by iterating over their fields.
		// Only process first-level anonymous fields.
		// Because the anonymous fields of the later level will be processed when the anonymous fields are parsed
		// In general, only fields at the current level are processed
		if fieldStruct.Anonymous {
			for j := 0; j < fieldStruct.Type.NumField(); j++ {
				fileName := GetFieldNameByTag(fieldStruct.Type.Field(j), op.tagName)

				if field.fieldName == TagIgnore || fieldStruct.Type.Field(j).Anonymous {
					continue
				}

				if _, ok := fields[fileName]; ok {
					continue
				}

				fields[fileName] = fieldInfo{
					fieldName:       fileName,
					fieldVal:        field.fieldVal.Field(j),
					belongAnonymous: field.fieldName,
				}
				// Associate the field with its embedded struct's name.
				anonymous[field.fieldName] = append(anonymous[field.fieldName], fileName)
			}
		}
	}
	return fields, anonymous
}

// get unexported field value
func getUnexportedField(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
