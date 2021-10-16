package collection

import (
	"reflect"
)

// HasAny returns true if the slice has more than one element
func HasAny(t interface{}) bool {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		return s.Len() > 0
	default:
		return false
	}
}
