package types

import (
	"reflect"
)

func IsNull(value interface{}) (isNull bool) {
	isNull = false
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) || v.IsNil()
	}
	return
}
