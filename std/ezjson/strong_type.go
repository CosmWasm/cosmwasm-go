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

	//getting tag
	Tag() string

	//setting values
	Set(value interface{}) error

	//setting attribute
	Attribute(key string, value bool)
}

type BaseName struct {
	tag         string
	realName    string
	omitempty   bool
	rust_option bool
}

func (b BaseName) Name() string {
	return b.realName
}

func (b BaseName) Tag() string {
	return b.tag
}

func (b *BaseName) Attribute(key string, value bool) {
	if key == OmitEmpty {
		b.omitempty = value
	} else if key == RustOption {
		b.rust_option = value
	}
}

func (b BaseName) IsOmitEmpty() bool {
	return b.omitempty
}

func (b BaseName) IsRustOption() bool {
	return b.rust_option
}

func (b BaseName) GetJsonName() string {
	if len(b.tag) > 0 {
		return b.tag
	}
	return b.realName
}

func quote(data string) string {
	return `"` + data + `"`
}

type BoolOpt struct {
	BaseName
	realValue bool
}

func (b BoolOpt) IsEmpty() bool {
	//boolean always return false
	return false
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
	if len(b.GetJsonName()) > 0 {
		result += quote(b.GetJsonName())
		result += ":"
	}
	if b.realValue == true {
		result += `"true"`
	} else {
		result += `"false"`
	}
	if isBrace == true {
		result += "}"
	}
	return result
}

func (b *BoolOpt) Set(value interface{}) error {
	b.realValue = value.(bool)
	return nil
}

///int int16 int32 int64 support
type IntOpt struct {
	BaseName
	realValue int64
}

func (i IntOpt) IsEmpty() bool {
	//int always return false
	return false
}

func (i IntOpt) Value() interface{} {
	return i.realValue
}

func (i IntOpt) Type() reflect.Kind {
	return reflect.Int
}

func (i IntOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(i.GetJsonName()) > 0 {
		result += quote(i.GetJsonName())
		result += ":"
	}
	result += strconv.FormatInt(i.realValue, 10)
	if isBrace == true {
		result += "}"
	}
	return result
}

func (i *IntOpt) Set(value interface{}) error {
	i.realValue = value.(int64)
	return nil
}

///int uint16 uint32 uint64 uintptr support
type UintOpt struct {
	BaseName
	realValue uint64
}

func (u UintOpt) IsEmpty() bool {
	//uint always return false
	return false
}

func (u UintOpt) Value() interface{} {
	return u.realValue
}

func (u UintOpt) Type() reflect.Kind {
	return reflect.Uint
}

func (u UintOpt) Encode(isBrace bool) string {
	result := ""
	if isBrace == true {
		result += "{"
	}
	if len(u.GetJsonName()) > 0 {
		result += quote(u.GetJsonName())
		result += ":"
	}
	result += strconv.FormatUint(u.realValue, 10)
	if isBrace == true {
		result += "}"
	}
	return result
}

func (u *UintOpt) Set(value interface{}) error {
	u.realValue = value.(uint64)
	return nil
}

///string support
type StringOpt struct {
	BaseName
	realValue string
}

func (s StringOpt) IsEmpty() bool {
	return len(s.realValue) == 0
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
	if len(s.GetJsonName()) > 0 {
		result += quote(s.GetJsonName())
		result += ":"
	}
	if s.IsRustOption() && len(s.realValue) == 0 {
		result += "null"
	} else {
		//todo need much strong conversion, add some check item, such as unrecognized symbols
		result += quote(s.realValue)
	}
	if isBrace == true {
		result += "}"
	}
	return result
}

func (s *StringOpt) Set(value interface{}) error {
	s.realValue = value.(string)
	return nil
}

//map support
type MapOpt struct {
}

//slice support
type SliceOpt struct {
	BaseName
	realValue []BaseOpt
}

func (s SliceOpt) IsEmpty() bool {
	return len(s.realValue) == 0
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
	if len(s.GetJsonName()) > 0 {
		result += quote(s.GetJsonName())
		result += ":"
	}
	if s.IsRustOption() && len(s.realValue) == 0 {
		result += "null"
	} else {
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
	}

	if isBrace == true {
		result += "}"
	}
	return result
}

func (s *SliceOpt) Set(value interface{}) error {
	s.realValue = decodeSlice(s.realName, s.tag, []byte(value.(string))).Value().([]BaseOpt)
	return nil
}

//struct support
type StructOpt struct {
	BaseName
	realValue []BaseOpt
}

func (s StructOpt) IsEmpty() bool {
	if len(s.realValue) == 0 {
		return true
	}
	for _, r := range s.realValue {
		if !r.IsEmpty() {
			return false //even only one of field was not empty, the structure was not empty too
		}
	}
	return true //return empty
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
	if len(s.GetJsonName()) > 0 {
		result += quote(s.GetJsonName())
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

func (s *StructOpt) Set(value interface{}) error {
	s.realValue = decodeStruct(s.tag, s.realName, []byte(value.(string))).Value().([]BaseOpt)
	return nil
}

//unsupport tyoe
type unsupportedOpt struct {
	BaseName
	realTypeName string
}

func (u unsupportedOpt) Init(data interface{}, name string) {
}

func (u unsupportedOpt) IsEmpty() bool {
	return false
}

func (u unsupportedOpt) Value() interface{} {
	return ""
}

func (u unsupportedOpt) Type() reflect.Kind {
	return reflect.Invalid
}

func (u unsupportedOpt) Encode(isBrace bool) string {
	return u.realTypeName
}

func (u unsupportedOpt) Set(value interface{}) error {
	return nil
}

//interface
func Generate(name, tag string, in interface{}, isDecoding bool) BaseOpt {
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
		return &BoolOpt{
			BaseName: BaseName{
				realName:    name,
				tag:         tag,
				omitempty:   false,
				rust_option: false,
			},
			realValue: ValueOf(ref).Bool(),
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &IntOpt{
			BaseName: BaseName{
				realName:    name,
				tag:         tag,
				omitempty:   false,
				rust_option: false,
			},
			realValue: ValueOf(ref).Int(),
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &UintOpt{
			BaseName: BaseName{
				realName:    name,
				tag:         tag,
				omitempty:   false,
				rust_option: false,
			},
			realValue: uint64(ValueOf(ref).Uint()),
		}
	case reflect.String:
		return &StringOpt{
			BaseName: BaseName{
				realName:    name,
				tag:         tag,
				omitempty:   false,
				rust_option: false,
			},
			realValue: ValueOf(ref).String(),
		}
	case reflect.Struct:
		p, e := prepare(ref, isDecoding)
		if e == nil {
			return &StructOpt{
				BaseName: BaseName{
					realName:    name,
					tag:         tag,
					omitempty:   false,
					rust_option: false,
				},
				realValue: p,
			}
		}
		return &unsupportedOpt{}
	case reflect.Slice, reflect.Array:
		p, e := prepare(ref, isDecoding)
		if e == nil {
			return &SliceOpt{
				BaseName: BaseName{
					realName:    name,
					tag:         tag,
					omitempty:   false,
					rust_option: false,
				},
				realValue: p,
			}
		}
		return &unsupportedOpt{}
	default:
		return &unsupportedOpt{
			BaseName: BaseName{
				realName:    name,
				tag:         tag,
				omitempty:   false,
				rust_option: false,
			},
			realTypeName: tinygo_typeof[kind],
		}
	}
}
