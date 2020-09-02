package ezjson

//for alpha version, we using jsonParser to support our unmarshal opt
//depending project : https://github.com/buger/jsonparser

import (
	"errors"
	"github.com/cosmwasm/jsonparser"
	"reflect"
)

func decodeJson(jsonstr []byte, out interface{}) error {
	if !IsPtr(out) {
		return errors.New("out put must be a pointer")
	}
	ref := reflect.ValueOf(out).Elem().Interface()
	opts, e := prepare(ref)
	if e != nil {
		return e
	}
	e = decoding(jsonstr, opts)
	if e != nil {
		return e
	}

	return assign(opts, out)
}

func lookup(name string, opts []BaseOpt) int {
	for idx, opt := range opts {
		if opt.Tag() == name || opt.Name() == name {
			return idx
		}
	}
	return -1
}

func queryRealValue(in []byte, dataType jsonparser.ValueType) interface{} {
	var e error
	var v interface{}
	switch dataType {
	case jsonparser.String:
		v, e = jsonparser.ParseString(in)
	case jsonparser.Number:
		v, e = jsonparser.ParseInt(in)
	case jsonparser.Boolean:
		v, e = jsonparser.ParseBoolean(in)
	default:
		v, e = jsonparser.ParseString(in)
	}
	if e == nil {
		return v
	}
	return nil
}

func decoding(jsonstr []byte, opts []BaseOpt) error {
	jsonparser.ObjectEach(jsonstr, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if opts != nil {
			idx := lookup(string(key), opts)
			if idx >= 0 {
				v := queryRealValue(value, dataType)
				if v != nil {
					opts[idx].Set(v)
				}
			}
		}
		return nil
	})
	return nil
}

func decodeStruct(name, tag string, jsonstr []byte) BaseOpt {
	var opts []BaseOpt
	jsonparser.ObjectEach(jsonstr, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if dataType == jsonparser.Array {
			opt := decodeSlice(string(key), string(key), value)
			if opt.Type() != reflect.Invalid {
				opts = append(opts, opt)
			}
			return nil
		} else if dataType == jsonparser.Object {
			opt := decodeStruct(string(key), string(key), value)
			if opt.Type() != reflect.Invalid {
				opts = append(opts, opt)
			}
			return nil
		}
		realValue := queryRealValue(value, dataType)
		if realValue == nil {
			return errors.New("Failed to query real value from data : " + string(key))
		}
		opt := Generate(string(key), string(key), realValue)
		if opt.Type() == reflect.Invalid {
			return nil //continue
		}
		opts = append(opts, opt)
		return nil
	})
	return &StructOpt{
		BaseName: BaseName{
			tag:      tag,
			realName: name,
		},
		realValue: opts,
	}
}

func decodeSlice(name, tag string, jsonstr []byte) BaseOpt {
	var opts []BaseOpt
	jsonparser.ArrayEach(jsonstr, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if dataType == jsonparser.Array {
			opt := decodeSlice("", "", value)
			if opt.Type() != reflect.Invalid {
				opts = append(opts, opt)
			}
			return
		} else if dataType == jsonparser.Object {
			opt := decodeStruct("", "", value)
			if opt.Type() != reflect.Invalid {
				opts = append(opts, opt)
			}
			return
		}
		realValue := queryRealValue(value, dataType)
		if realValue == nil {
			return //continue
		}
		opt := Generate("", "", realValue)
		if opt.Type() == reflect.Invalid {
			return //continue
		}
		opts = append(opts, opt)
	})
	return &SliceOpt{
		BaseName: BaseName{
			tag:      tag,
			realName: name,
		},
		realValue: opts,
	}
}

func doAssign(opts []BaseOpt, vals reflect.Value, tps reflect.Type) error {
	if tps.Kind() == reflect.Slice || tps.Kind() == reflect.Array {
		Log("Process Slice")
		if len(opts) <= 0 {
			//if none, skip it
			return nil
		}
		//must equal with original type
		var boolSlice []bool
		var intSlice []int
		var int8Slice []int8
		var int16Slice []int16
		var int32Slice []int32
		var int64Slice []int64
		var uintSlice []uint
		var uint8Slice []uint8
		var uint16Slice []uint16
		var uint32Slice []uint32
		var uint64Slice []uint64
		var uintptrSlice []uintptr
		var stringSlice []string

		for _, opt := range opts {
			switch tps.Elem().Kind() {
			case reflect.Bool:
				boolSlice = append(boolSlice, opt.Value().(bool))
			case reflect.Int:
				intSlice = append(intSlice, int(opt.Value().(int64)))
			case reflect.Int8:
				int8Slice = append(int8Slice, int8(opt.Value().(int64)))
			case reflect.Int16:
				int16Slice = append(int16Slice, int16(opt.Value().(int64)))
			case reflect.Int32:
				int32Slice = append(int32Slice, int32(opt.Value().(int64)))
			case reflect.Int64:
				int64Slice = append(int64Slice, opt.Value().(int64))
			case reflect.Uint:
				uintSlice = append(uintSlice, uint(opt.Value().(uint64)))
			case reflect.Uint8:
				uint8Slice = append(uint8Slice, uint8(opt.Value().(uint64)))
			case reflect.Uint16:
				uint16Slice = append(uint16Slice, uint16(opt.Value().(uint64)))
			case reflect.Uint32:
				uint32Slice = append(uint32Slice, uint32(opt.Value().(uint64)))
			case reflect.Uint64:
				uint64Slice = append(uint64Slice, opt.Value().(uint64))
			case reflect.Uintptr:
				uintptrSlice = append(uintptrSlice, uintptr(opt.Value().(uint64)))
			case reflect.String:
				stringSlice = append(stringSlice, opt.Value().(string))
			case reflect.Struct:
				doAssign(opt.Value().([]BaseOpt), vals, tps.Elem())
			case reflect.Slice, reflect.Array:
				if opt.Type() != reflect.Slice && opt.Type() != reflect.Array {
					if opt.IsEmpty() {
						//if value is empty, we skip it
						continue
					}
					if opt.Type() == reflect.String {
						//if tart value is string, we can try to translate it, otherwise, return error
						trv := []byte(opt.Value().(string))
						trvopt := Generate(opt.Name(), opt.Tag(), trv)
						if trvopt.Type() == reflect.Invalid {
							return errors.New("decode failed, wrong param passed, found" + tinygo_typeof[opt.Type()] + " excepted Slice or Array")
						}
						doAssign(trvopt.Value().([]BaseOpt), vals, tps.Elem()) //make it to slice
					} else {
						continue
						//	return errors.New("decode failed, wrong param passed, found" + tinygo_typeof[opt.Type()] + " excepted Slice or Array")
					}
				} else {
					doAssign(opt.Value().([]BaseOpt), vals, tps.Elem())
				}
			default:
				continue
			}
		}

		switch tps.Elem().Kind() {
		case reflect.Bool:
			vals.Set(ValueOf(boolSlice))
		case reflect.Int:
			vals.Set(ValueOf(intSlice))
		case reflect.Int8:
			vals.Set(ValueOf(int8Slice))
		case reflect.Int16:
			vals.Set(ValueOf(int16Slice))
		case reflect.Int32:
			vals.Set(ValueOf(int32Slice))
		case reflect.Int64:
			vals.Set(ValueOf(int64Slice))
		case reflect.Uint:
			vals.Set(ValueOf(uintSlice))
		case reflect.Uint8:
			vals.Set(ValueOf(uint8Slice))
		case reflect.Uint16:
			vals.Set(ValueOf(uint16Slice))
		case reflect.Uint32:
			vals.Set(ValueOf(uint32Slice))
		case reflect.Uint64:
			vals.Set(ValueOf(uint64Slice))
		case reflect.Uintptr:
			vals.Set(ValueOf(uintptrSlice))
		case reflect.String:
			vals.Set(ValueOf(stringSlice))
		}

		return nil
	}
	FieldLen := vals.NumField()
	for i := 0; i < FieldLen; i++ {
		tp := tps.Field(i)
		realName, _ := getTag(string(tp.Tag))
		if len(realName) <= 0 {
			realName = tps.Field(i).Name
		}
		idx := lookup(realName, opts)
		if idx >= 0 {
			opt := opts[idx]
			val := vals.Field(i)
			switch tp.Type.Kind() {
			case reflect.Bool:
				val.SetBool(opt.Value().(bool))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val.SetInt(opt.Value().(int64))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				var value uint64
				if opt.Type() == reflect.Int ||
					opt.Type() == reflect.Int8 ||
					opt.Type() == reflect.Int16 ||
					opt.Type() == reflect.Int32 ||
					opt.Type() == reflect.Int64 {
					value = uint64(opt.Value().(int64)) //in decoding, only int type
				} else {
					value = opt.Value().(uint64) //in decoding, only int type
				}
				val.SetUint(value)
			case reflect.String:
				val.SetString(opt.Value().(string))
			case reflect.Struct:
				doAssign(opt.Value().([]BaseOpt), val, tps.Field(i).Type)
			case reflect.Slice, reflect.Array:
				if opt.Type() != reflect.Slice && opt.Type() != reflect.Array {
					if opt.IsEmpty() {
						//if value is empty, we skip it
						continue
					}
					if opt.Type() == reflect.String {
						//if tart value is string, we can try to translate it, otherwise, return error
						trv := []byte(opt.Value().(string))
						trvopt := Generate(opt.Name(), opt.Tag(), trv)
						if trvopt.Type() == reflect.Invalid {
							continue
							//	return errors.New("decode failed, wrong param passed, found" + tinygo_typeof[opt.Type()] + " excepted Slice or Array")
						}
						doAssign(trvopt.Value().([]BaseOpt), val, tps.Field(i).Type) //make it to slice
					} else {
						return errors.New("decode failed, wrong param passed, found" + tinygo_typeof[opt.Type()] + " excepted Slice or Array")
					}
				} else {
					doAssign(opt.Value().([]BaseOpt), val, tps.Field(i).Type)
				}
			default:
				continue
			}
		}
	}
	return nil
}

func assign(opts []BaseOpt, out interface{}) error {
	tps := reflect.TypeOf(out)
	vals := reflect.ValueOf(out)
	if tps.Kind() == reflect.Ptr {
		tps = tps.Elem()
	}
	if vals.Kind() == reflect.Ptr {
		vals = vals.Elem()
	}
	return doAssign(opts, vals, tps)
}
