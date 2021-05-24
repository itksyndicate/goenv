package goenv

import "reflect"

func isPtrType(value reflect.Type) bool {
	return value.Kind() == reflect.Ptr
}

func isStructType(value reflect.Type) bool {
	return value.Kind() == reflect.Struct
}

func isStructPtrType(value reflect.Type) bool {
	return isPtrType(value) && isStructType(value.Elem())
}
