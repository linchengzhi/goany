package goany

import (
	"reflect"
	"strings"
)

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

// GetFieldNameByTag returns the field name by tag.
// if tag is gorm, get column name
// if tag name is nil, get field name
func GetFieldNameByTag(field reflect.StructField, tag string) string {
	tagValue := field.Tag.Get(tag)
	if tagValue == "" {
		return field.Name
	}

	if tag == "gorm" {
		columnPrefix := "column:"
		start := strings.Index(tagValue, columnPrefix)
		if start != -1 {
			tagValue = tagValue[start+len(columnPrefix):]
		}
	}

	if i := strings.Index(tagValue, ","); i != -1 {
		tagValue = tagValue[:i]
	}
	return tagValue
}
