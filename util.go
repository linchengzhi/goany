package goany

import (
	"reflect"
	"strings"
)

// CheckInIsNil checks if the input is nil or if it's a nil pointer.
func CheckInIsNil(in interface{}) bool {
	if in == nil {
		return true
	}
	val := reflect.ValueOf(in)
	return val.Kind() == reflect.Ptr && val.IsNil()
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
