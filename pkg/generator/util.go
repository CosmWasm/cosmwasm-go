package generator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"reflect"
)

// reflectTypeToGoType converts a reflect.Type to go type.
// TODO(fdymylja): support type definition over primitive types [type Custom string] and type aliasing [type Something = string]
func reflectTypeToGoType(typ reflect.Type) string {
	kind := typ.Kind()
	if kind == reflect.Slice && typ.Elem().Kind() != reflect.Uint8 {
		panic(fmt.Errorf("slice types are supported only for bytes keys"))
	}
	switch kind {
	case reflect.Bool:
		return "bool"
	case reflect.Int8:
		return "int8"
	case reflect.Int16:
		return "int16"
	case reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Uint8:
		return "uint8"
	case reflect.Uint16:
		return "uint16"
	case reflect.Uint32:
		return "uint32"
	case reflect.Uint64:
		return "uint64"
	case reflect.Slice: // only for bytes
		return "[]byte"
	case reflect.String:
		return "string"
	default:
		panic(fmt.Errorf("unsupported key kind: %s", kind.String()))
	}
}

func lowerCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}
