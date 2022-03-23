package collection

import (
	"reflect"

	"github.com/google/uuid"
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

// Contains returns true if the slice contains the element
func Contains(t interface{}, element interface{}) bool {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == element {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func UUIDDiffs(prevIDs []uuid.UUID, currentIDs []uuid.UUID) []uuid.UUID {
	var diffs []uuid.UUID
	for _, currentID := range currentIDs {
		if !Contains(prevIDs, currentID) {
			diffs = append(diffs, currentID)
		}
	}
	return diffs
}
