package ezjson

import "reflect"

//give a determinacy reflect support between go and tinyGo
//that make coding process clearly

func TypeOf(in interface{}) reflect.Type {
	return reflect.TypeOf(in)
}

func ValueOf(in interface{}) reflect.Value {
	return reflect.ValueOf(in)
}

func IsPtr(in interface{}) bool {
	return TypeOf(in).Kind() == reflect.Ptr
}

func IsUnsafePointer(in interface{}) bool {
	return TypeOf(in).Kind() == reflect.UnsafePointer
}

func AllElem(in interface{}) reflect.Value {
	if IsPtr(in) {
		return ValueOf(in).Elem()
	}
	return ValueOf(in)
}
