package ezjson

import (
	"reflect"
	"strconv"
)

type BaseOpt interface {
	//checking is value empty
	IsEmpty() bool

	//get encoded result
	Encode(isBrace bool) string

	//getting Value name
	Name() string

	//getting value
	Value() interface{}

	//getting type
	Type() reflect.Kind
}

func quote(data string) string {
	return `"` + data + `"`
}

type BoolOpt struct {
	realValue bool
	realName  string
}

func (b BoolOpt) IsEmpty() bool {
	//boolean always return false
	return false
}

func (b BoolOpt) Name() string {
	return b.realName
}

func (b BoolOpt) Value() interface{} {
	return b.realValue
}

func (b BoolOpt) Type() reflect.Kind {
	return reflect.Bool
}

func (b BoolOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(b.realName) > 0 {
		result += quote(b.realName)
		result += ":"
	}
	if b.realValue == true {
		result += "1"
	} else {
		result += "0"
	}
	if isBrace == true {
		result += "}"
	}
	return result
}

///int int16 int32 int64 support
type IntOpt struct {
	realValue int64
	realName  string
}

func (i IntOpt) IsEmpty() bool {
	//boolean always return false
	return false
}

func (i IntOpt) Name() string {
	return i.realName
}

func (i IntOpt) Value() interface{} {
	return i.realValue
}

func (i IntOpt) Type() reflect.Kind {
	return reflect.Bool
}

func (i IntOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(i.realName) > 0 {
		result += quote(i.realName)
		result += ":"
	}
	result += strconv.FormatInt(i.realValue, 10)
	if isBrace == true {
		result += "}"
	}
	return result
}

///int uint16 uint32 uint64 uintptr support
type UintOpt struct {
	realValue uint64
	realName  string
}

func (u UintOpt) IsEmpty() bool {
	//boolean always return false
	return false
}

func (u UintOpt) Name() string {
	return u.realName
}

func (u UintOpt) Value() interface{} {
	return u.realValue
}

func (u UintOpt) Type() reflect.Kind {
	return reflect.Bool
}

func (u UintOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(u.realName) > 0 {
		result += quote(u.realName)
		result += ":"
	}
	result += strconv.FormatUint(u.realValue, 10)
	if isBrace == true {
		result += "}"
	}
	return result
}

///string support
type StringOpt struct {
	realValue string
	realName  string
}

func (s StringOpt) IsEmpty() bool {
	//boolean always return false
	return false
}

func (s StringOpt) Name() string {
	return s.realName
}

func (s StringOpt) Value() interface{} {
	return s.realValue
}

func (s StringOpt) Type() reflect.Kind {
	return reflect.String
}

func (s StringOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(s.realName) > 0 {
		result += quote(s.realName)
		result += ":"
	}
	//todo need much strong conversion, add some check item, such as unrecognized symbols
	result += quote(s.realValue)
	if isBrace == true {
		result += "}"
	}
	return result
}

//map support
type MapOpt struct {
}

//slice support
type SliceOpt struct {
	realValue []BaseOpt
	realName  string
}

func (s SliceOpt) IsEmpty() bool {
	return false
}

func (s SliceOpt) Name() string {
	return s.realName
}

func (s SliceOpt) Value() interface{} {
	return s.realValue
}

func (s SliceOpt) Type() reflect.Kind {
	return reflect.Slice
}

func (s SliceOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(s.realName) > 0 {
		result += quote(s.realName)
		result += ":"
	}
	jsonstr := "["
	first := true
	for _, opt := range s.realValue {
		if first == false {
			jsonstr += ","
		}
		jsonstr += opt.Encode(false)
		if first == true {
			first = false
		}
	}
	jsonstr += "]"
	result += jsonstr
	if isBrace == true {
		result += "}"
	}
	return result
}

//struct support
type StructOpt struct {
	realValue []BaseOpt
	realName  string
}

func (s StructOpt) IsEmpty() bool {
	return false
}

func (s StructOpt) Name() string {
	return s.realName
}

func (s StructOpt) Value() interface{} {
	return s.realValue
}

func (s StructOpt) Type() reflect.Kind {
	return reflect.Struct
}

func (s StructOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(s.realName) > 0 {
		result += quote(s.realName)
		result += ":"
	}
	jsonstr := "{"
	first := true
	for _, opt := range s.realValue {
		if first == false {
			jsonstr += ","
		}
		jsonstr += opt.Encode(false)
		if first == true {
			first = false
		}
	}
	jsonstr += "}"
	result += jsonstr
	if isBrace == true {
		result += "}"
	}
	return result
}

//unsupport tyoe
type unsupportedOpt struct {
	realName     string
	realTypeName string
}

func (u unsupportedOpt) Init(data interface{}, name string) {
}

func (s unsupportedOpt) IsEmpty() bool {
	return false
}

func (s unsupportedOpt) Name() string {
	return s.realName
}

func (s unsupportedOpt) Value() interface{} {
	return ""
}

func (s unsupportedOpt) Type() reflect.Kind {
	return reflect.Invalid
}

func (s unsupportedOpt) Encode(isBrace bool) string {
	return s.realTypeName
}

//interface
func Generate(name string, in interface{}) BaseOpt {
	ref := in
	if IsPtr(in) {
		//under go, it works well, in tinyGo, Typeof(ref).Kind will return Invalid
		//so we denied to using ptr in tinyGo
		//ref = (ValueOf(in)).Elem().Interface()
		//todo must support later
	}
	kind := TypeOf(ref).Kind()
	switch kind {
	case reflect.Bool:
		return BoolOpt{
			realValue: ValueOf(ref).Bool(),
			realName:  name,
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return IntOpt{
			realValue: ValueOf(ref).Int(),
			realName:  name,
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return UintOpt{
			realValue: ValueOf(ref).Uint(),
			realName:  name,
		}
	case reflect.String:
		return StringOpt{
			realValue: ValueOf(ref).String(),
			realName:  name,
		}
	case reflect.Struct:
		p, e := prepare(ref)
		if e == nil {
			return StructOpt{
				realValue: p,
				realName:  name,
			}
		}
		return unsupportedOpt{}
	case reflect.Slice, reflect.Array:
		p, e := prepare(ref)
		if e == nil {
			return SliceOpt{
				realValue: p,
				realName:  name,
			}
		}
		return unsupportedOpt{}
	case reflect.Float32, reflect.Float64:
		return unsupportedOpt{
			realName:     name,
			realTypeName: tinygo_typeof[kind],
		}
	default:
		return unsupportedOpt{
			realName:     name,
			realTypeName: tinygo_typeof[kind],
		}
	}
}
