package goany

import (
	"reflect"
	"unsafe"
)

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
func deepInFields(in interface{}, anonName string, op Options) ([]fieldInfo, error) {
	_, inVal := ReflectTypeValue(in)
	if inVal.Kind() != reflect.Struct {
		return []fieldInfo{}, nil
	}
	fields := make([]fieldInfo, 0, inVal.NumField())
	for i := 0; i < inVal.NumField(); i++ {

		field := new(fieldInfo)
		field.fieldStruct = inVal.Type().Field(i)
		field.fieldVal = inVal.Field(i)
		field.fieldName = GetFieldNameByTag(field.fieldStruct, op.tagName)

		if !field.canUse(op) {
			continue
		}
		// If the field is an anonName struct, use unsafe to get val
		if field.fieldStruct.PkgPath != "" {
			if !field.fieldVal.CanAddr() {
				return []fieldInfo{}, ErrInNotPtr
			}
			field.fieldVal = getUnexportedField(field.fieldVal)
		}

		if field.fieldVal.Kind() == reflect.Struct && op.exportedUnExported {
			field.fieldVal = reflect.NewAt(field.fieldStruct.Type, unsafe.Pointer(field.fieldVal.UnsafeAddr()))
		}
		// Retrieve unexported field values if necessary.
		field.belongAnonymous = anonName
		fields = append(fields, *field)
		if field.fieldStruct.Anonymous && anonName == "" { //anonName field, only first level
			anonFields, _ := deepInFields(field.fieldVal.Interface(), field.fieldName, op)
			fields = append(fields, anonFields...)
		}
	}
	return fields, nil
}

// deepOutFields creates a mapping of field names to fieldInfo structs and a mapping of
// embedded struct names to their respective fields' names for an output struct. It processes
// each field in the struct, including those in embedded structs, and handles unexported fields
// if the exportedLower option is set.
func deepOutFields(outVal reflect.Value, anonName string, op Options) (map[string]fieldInfo, map[string][]string) {
	outVal, _ = indirectValue(outVal)
	if outVal.Kind() != reflect.Struct {
		return map[string]fieldInfo{}, map[string][]string{}
	}
	fields := make(map[string]fieldInfo)
	anonymous := make(map[string][]string)
	for i := 0; i < outVal.NumField(); i++ {
		field := new(fieldInfo)
		field.fieldStruct = outVal.Type().Field(i)
		field.fieldVal = outVal.Field(i)
		field.fieldName = GetFieldNameByTag(field.fieldStruct, op.tagName)

		if !field.canUse(op) {
			continue
		}

		if field.fieldVal.Kind() == reflect.Ptr && field.fieldVal.IsNil() { //if field is nil, init
			if field.fieldStruct.PkgPath == "" {
				field.fieldVal.Set(reflect.New(field.fieldStruct.Type.Elem()))
			} else {
				field.fieldVal = getUnexportedField(field.fieldVal)
			}
		}

		// Add or update the fieldInfo in the fields mapping.
		if f, ok := fields[field.fieldName]; !ok || f.belongAnonymous != "" {
			fields[field.fieldName] = *field
		}

		// Handle embedded (anonymous) struct fields by iterating over their fields.
		// Only process first-level anonymous fields.
		// Because the anonymous fields of the later level will be processed when the anonymous fields are parsed
		// In general, only fields at the current level are processed
		if field.fieldStruct.Anonymous && anonName == "" {
			// If the field is an anonymous struct, use unsafe to get val
			reflectType, _ := ReflectTypeValue(field.fieldStruct.Type)
			if reflectType.Kind() == reflect.Struct {
				anonFields, annoFieldNames := deepOutFields(field.fieldVal, field.fieldName, op)
				for k, v := range anonFields {
					if _, ok := fields[k]; !ok {
						fields[k] = v
					}
				}
				for k, v := range annoFieldNames {
					if _, ok := anonymous[k]; !ok {
						anonymous[k] = v
					}
				}
			}
		} else if anonName != "" {
			anonymous[anonName] = append(anonymous[anonName], field.fieldName)
		}
	}
	return fields, anonymous
}

// canUse returns whether the field can be used
func (field *fieldInfo) canUse(op Options) bool {
	if field.fieldName == TagIgnore {
		return false
	}
	if field.fieldStruct.PkgPath != "" && !op.exportedUnExported {
		return false
	}
	return true
}
