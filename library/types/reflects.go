package types

import (
	"reflect"
)

func NewVarByInterface(obj interface{}) interface{} {
	ptr := reflect.New(reflect.TypeOf(obj))
	return ptr.Elem().Interface()
}

func NewVarAddrByInterface(addr interface{}) interface{} {
	var (
		typeOf = reflect.TypeOf(addr).Elem()
		ptr    = reflect.New(typeOf)
		v      = ptr.Elem().Interface()
	)
	return &v
}
