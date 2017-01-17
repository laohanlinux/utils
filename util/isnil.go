package util 

import (
	"reflect"
)

func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}

	// if object type is nil, value will be zero value
	value := reflect.ValueOf(object)
	kind := value.Kind()
	
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}
