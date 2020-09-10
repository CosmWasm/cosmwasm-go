package ezjson

import (
	"errors"
	"reflect"
)

func package_by_type(in interface{}, data string) string {
	kind := TypeOf(in).Kind()
	switch kind {
	case reflect.Struct:
		return "{" + data + "}"
	case reflect.Slice, reflect.Array:
		return "[" + data + "]"
	case reflect.Map:
		return "{" + data + "}"
	default:
		return data
	}
}

func encode2json(in interface{}) ([]byte, error) {
	opts, err := prepare(in, false)
	if err != nil {
		return nil, err
	}
	jsonstr := ""
	first := true
	for _, opt := range opts {
		if first == false {
			jsonstr += ","
		}
		jsonstr += opt.Encode(false)
		if first == true {
			first = false
		}
	}
	result := package_by_type(in, jsonstr)
	return []byte(result), nil
}

func prepare(in interface{}, isDecoding bool) ([]BaseOpt, error) {
	opts := make([]BaseOpt, 0)
	t := reflect.TypeOf(in)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		vals := reflect.ValueOf(in)

		for i := 0; i < t.NumField(); i++ {
			field := vals.Field(i)
			if field.CanInterface() {
				tag, isOmit := getTag(string(t.Field(i).Tag))
				name := t.Field(i).Name
				vi := field.Interface()
				opt := Generate(name, tag, vi, isDecoding)
				if opt.Type() == reflect.Invalid {
					return nil, errors.New("Error : Invalid conversion type :[" + t.Field(i).Name + "] -- [" + opt.Encode(false) + "]")
				}
				if !isDecoding && isOmit && opt.IsEmpty() {
					continue //skip by omitempty key word if target value is empty
				}
				opts = append(opts, opt)
			} else {
				//is ptr?
				//never happened, in tinygo, CanInterface always return true
			}
		}
	} else if t.Kind() == reflect.Slice {
		vals := reflect.ValueOf(in)
		if vals.Kind() == reflect.Ptr {
			vals = vals.Elem()
		}
		var opt BaseOpt
		for i := 0; i < vals.Len(); i++ {
			opt = Generate("", "", vals.Index(i).Interface(), isDecoding)
			if opt.Type() == reflect.Invalid {
				return nil, errors.New("Error : Invalid conversion slice type")
			}
			opts = append(opts, opt)
		}

	}
	return opts, nil
}
