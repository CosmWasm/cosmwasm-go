// Package vjson is a small, minimal alternative to encoding/json which has no dependencies
// and works with Tingyo.  The Marshal and Unmarshal methods work like encoding/json, but
// structs are not supported, only primitives plus map[string]interface{} and []interface{}.
package vjson

import (
	"bytes"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/poc/ext/jsonparser"
	"io"
	"reflect"
	"strconv"
)

// NOTE: various bits have been borrowed from encoding/json

// bool, for JSON booleans
// string, for JSON strings
// []interface{}, for JSON arrays
// map[string]interface{}, for JSON objects
// nil for JSON null
// RawMessage - just use Marshaler

// if someone asks to ask into an int type that should still work

func marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := marshalTo(&buf, v)
	return buf.Bytes(), err
}

func marshalEx(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := marshalToEx(&buf, v)
	return buf.Bytes(), err
}

func marshalToEx(w io.Writer, vin interface{}) (err error){
	bb := make([]byte, 0, 128) // hopefully stack alloc

	typeof := reflect.TypeOf(vin)
	if typeof == reflect.TypeOf(bool(true)) {
		val := reflect.ValueOf(vin)
		bb = strconv.AppendBool(bb, val.Bool())
	}else if typeof == reflect.TypeOf(int(0)) ||
		typeof == reflect.TypeOf(int8(0)) ||
		typeof == reflect.TypeOf(int16(0))||
		typeof == reflect.TypeOf(int32(0))||
		typeof == reflect.TypeOf(int64(0)){
			val := reflect.ValueOf(vin)
			bb = strconv.AppendInt(bb, val.Int(),10)
	}else if typeof == reflect.TypeOf(uint(0))||
		typeof == reflect.TypeOf(uint8(0)) ||
		typeof == reflect.TypeOf(uint16(0)) ||
		typeof == reflect.TypeOf(uint32(0)) ||
		typeof == reflect.TypeOf(uint64(0)) {
			val := reflect.ValueOf(vin)
			bb = strconv.AppendUint(bb, val.Uint(),10)
	}else if typeof == reflect.TypeOf(string("")){
		val := reflect.ValueOf(vin)
		return encodeString(w, val.String(), false)
	}else{		//default to interface{}
		t := reflect.TypeOf(vin)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		first := true
		if t.Kind()== reflect.Struct {
			w.Write([]byte(`{`))
			vals := reflect.ValueOf(vin)

			if vals.Kind() == reflect.Ptr {
				vals = vals.Elem()
			}
			for i := 0; i < t.NumField(); i++ {
				if vals.Field(i).CanInterface(){
					if !first {
						w.Write([]byte(`,`))
					}
					data := vals.Field(i).Interface()
					w.Write([]byte(`"`))
					w.Write([]byte(t.Field(i).Name))
					w.Write([]byte(`":`))
					vkind := vals.Field(i).Kind()
					if (vkind == reflect.Ptr || vkind == reflect.UnsafePointer) && vals.Field(i).IsNil() {
						w.Write([]byte(`null`))
					}else{
						err := marshalToEx(w, data)
						if err != nil {
							return err
						}
					}
					first = false
				}
			}
			_, err := w.Write([]byte(`}`))
			return err
		}
	}
	if len(bb) == 0 {
		return errors.New("unexpected zero length buffer")
	}
	_, err = w.Write(bb)
	return err

}

func marshalTo(w io.Writer, vin interface{}) (err error) {

	bb := make([]byte, 0, 64) // hopefully stack alloc

	// TODO: check for Marshaler

	// TODO: check for MarshalerTo

	if vin == nil {
		_, err = w.Write([]byte("null"))
		return err
	}

	switch v := vin.(type) {

	case bool:
		bb = strconv.AppendBool(bb, v)

	case int:
		bb = strconv.AppendInt(bb, int64(v), 10)
	case int8:
		bb = strconv.AppendInt(bb, int64(v), 10)
	case int16:
		bb = strconv.AppendInt(bb, int64(v), 10)
	case int32:
		bb = strconv.AppendInt(bb, int64(v), 10)
	case int64:
		bb = strconv.AppendInt(bb, v, 10)
	case uint:
		bb = strconv.AppendUint(bb, uint64(v), 10)
	case uint8:
		bb = strconv.AppendUint(bb, uint64(v), 10)
	case uint16:
		bb = strconv.AppendUint(bb, uint64(v), 10)
	case uint32:
		bb = strconv.AppendUint(bb, uint64(v), 10)
	case uint64:
		bb = strconv.AppendUint(bb, v, 10)

	case string:
		return encodeString(w, v, false)

	// case []byte: // TODO: this is wrong - byte slice should get base64 encoded
	// 	return encodeStringBytes(w, v, false)

	case []interface{}:
		if v == nil {
			bb = append(bb, `null`...)
			break
		}
		w.Write([]byte(`[`))
		first := true
		for i := range v {
			if !first {
				w.Write([]byte(`,`))
			}
			first = false
			err := marshalTo(w, v[i])
			if err != nil {
				return err
			}
		}
		_, err := w.Write([]byte(`]`))
		return err

	case map[string]interface{}:
		if v == nil {
			bb = append(bb, `null`...)
			break
		}
		w.Write([]byte(`{`))
		first := true
		for k, el := range v {
			if !first {
				w.Write([]byte(`,`))
			}
			first = false
			encodeString(w, k, false)
			w.Write([]byte(`:`))
			err := marshalTo(w, el)
			if err != nil {
				return err
			}
		}
		_, err := w.Write([]byte(`}`))
		return err
	case interface{}:
		{
			t := reflect.TypeOf(vin)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			first := true
			if t.Kind()== reflect.Struct {
				w.Write([]byte(`{`))
				vals := reflect.ValueOf(vin)

				if vals.Kind() == reflect.Ptr {
					vals = vals.Elem()
				}
				for i := 0; i < t.NumField(); i++ {
					if vals.Field(i).CanInterface(){
						if !first {
							w.Write([]byte(`,`))
						}
						data := vals.Field(i).Interface()
						w.Write([]byte(`"`))
						w.Write([]byte(t.Field(i).Name))
						w.Write([]byte(`":`))
						vkind := vals.Field(i).Kind()
						if (vkind == reflect.Ptr || vkind == reflect.UnsafePointer) && vals.Field(i).IsNil() {
							w.Write([]byte(`null`))
						}else{
							err := marshalTo(w, data)
							if err != nil {
								return err
							}
						}
						first = false
					}
				}
				_, err := w.Write([]byte(`}`))
				return err
			}

		}

	// TODO: pointer cases

	default:
		return errors.New("vjson.marshalTo error unknown type")
	}

	if len(bb) == 0 {
		panic("unexpected zero length buffer") // should never happen
	}

	_, err = w.Write(bb)
	return err
}

func unmarshal(data []byte, v interface{}) error {
	return unmarshalFrom(data,bytes.NewReader(data), v)
}

func unmarshalStruct(data []byte, vin interface{}) error{

	tps := reflect.TypeOf(vin)
	vals := reflect.ValueOf(vin)
	if tps.Kind() == reflect.Ptr{
		tps = tps.Elem()
	}
	if vals.Kind() == reflect.Ptr {
		vals = vals.Elem()
	}
	for i := 0; i < vals.NumField(); i++ {
		tp := tps.Field(i).Type
		valName := tps.Field(i).Name
		if tp == reflect.TypeOf(int(0)) {
			rint,e := jsonparser.GetInt(data,valName)
			if e == nil {
				vals.Field(i).SetInt(rint)
			}else {
				return errors.New("Can not find ["+valName + "], error: " + e.Error())
			}
		}else if tp == reflect.TypeOf(string("0")){
			rstr,e := jsonparser.GetString(data,valName)
			if e == nil {
				vals.Field(i).SetString(rstr)
			}else {
				return errors.New("Can not find ["+valName + "], error: " + e.Error())
			}
		}else {
			return errors.New("unsupported struct type")
		}
	}
	return nil
}

func unmarshalFrom(data []byte,r unreader, vin interface{}) error {

	// read the next token, whatever it is
	tok, err := readToken(r)
	if err != nil {
		return err
	}

	return unmarshalNext(data,r, tok, vin)
}

func unmarshalNext(data []byte, r unreader, tok Token, vin interface{}) error {

	tokDelim, _ := tok.(Delim)

	// and type switch to determine how to handle it
	switch v := vin.(type) {

	case interface{}:
		{
			tps := reflect.TypeOf(vin)
			vals := reflect.ValueOf(vin)
			if tps.Kind() == reflect.Ptr{
				tps = tps.Elem()
			}
			if vals.Kind() == reflect.Ptr {
				vals = vals.Elem()
			}
			for i := 0; i < vals.NumField(); i++ {
				tp := tps.Field(i).Type
				valName := tps.Field(i).Name
				if tp == reflect.TypeOf(int(0)) {
					rint,e := jsonparser.GetInt(data,valName)
					if e == nil {
						vals.Field(i).SetInt(rint)
					}
				}else if tp == reflect.TypeOf(string("0")){
					rstr,e := jsonparser.GetString(data,valName)
					if e == nil {
						vals.Field(i).SetString(rstr)
					}
				}
			}
		}
	case *bool:
		if tokv, ok := tok.(bool); ok {
			*v = tokv
		} else {
			return errors.New("vjson.unmarshalNext unable to scan")
		}


	// case *int:

	// case *int8:

	// case *int16:

	// case *int32:

	// case *int64:

	// case *uint:

	// case *uint8:

	// case *uint16:

	// case *uint32:

	// case *uint64:

	case *string:
		if tokv, ok := tok.(string); ok {
			*v = tokv
		} else {
			return errors.New("vjson.unmarshalNext unable to scan")
		}

	// case *[]byte: // hm, should be base64 encoded
	// 	if tokv, ok := tok.(string); ok {
	// 		*v = []byte(tokv)
	// 	} else {
	// 		return fmt.Errorf("vjson.unmarshalNext unable to scan %#v into %T", tok, v)
	// 	}

	case *[]interface{}:

		// make sure we have an array start
		if tokDelim != Delim('[') {
			return errors.New("vjson.unmarshalNext unable to scan")
		}

		sliceV := make([]interface{}, 0, 4)

		for {
			nextTok, err := readToken(r)
			if err != nil {
				return err
			}
			nextTokDelim, _ := nextTok.(Delim)
			if nextTokDelim == Delim(']') {
				break // end array
			}
			elV := newDefaultForToken(nextTok)
			err = unmarshalNext(data,r, nextTok, elV)
			if err != nil {
				return err
			}
			sliceV = append(sliceV, deref(elV))
		}

		*v = sliceV

	case *map[string]interface{}:

		// make sure we have an object start
		if tokDelim != Delim('{') {
			return errors.New("vjson.unmarshalNext unable to scan")
		}

		mapV := make(map[string]interface{}, 4)

		for {
			// read object key (must be string)
			keyTok, err := readToken(r)
			if err != nil {
				return err
			}
			keyTokDelim, _ := keyTok.(Delim)
			if keyTokDelim == Delim('}') {
				break // end object
			}

			nextTok, err := readToken(r)
			if err != nil {
				return err
			}

			elV := newDefaultForToken(nextTok)
			err = unmarshalNext(data,r, nextTok, elV)
			if err != nil {
				return err
			}

			keyTokStr, ok := keyTok.(string)
			if !ok {
				return errors.New("unexpected non-string object key token")
			}
			mapV[keyTokStr] = deref(elV)

		}

		*v = mapV

	default:
		return errors.New("vjson.unmarshalNext error unknown type ")

	}

	return nil
}

// newDefaultForToken will return a pointer to the appropriate type based on a JSON token.
// Used when scanning into an interface{} and we need to infer the Go type from the JSON input.
// A nil Token will return nil.
func newDefaultForToken(tok Token) interface{} {

	if tok == nil {
		return new(interface{})
	}

	switch tok.(type) {
	case bool:
		return new(bool)
	case Number:
		return new(int64)
	case string:
		return new(string)
	}

	tokDelim, _ := tok.(Delim)
	if tokDelim == Delim('[') {
		return new([]interface{})
	} else if tokDelim == Delim('{') {
		return new(map[string]interface{})
	}

	panic("newDefaultForToken unexpected token")
}

// deref will strip the pointer off of the value returned by newDefaultForToken
func deref(vin interface{}) interface{} {

	// if vin == nil {
	// 	return nil
	// }

	switch v := vin.(type) {
	case *interface{}:
		if *v == nil {
			return nil
		} else {
			panic("deref: *interface{} should have been nil but got")
		}
	case *bool:
		return *v
	case *string:
		return *v
	case *[]interface{}:
		return *v
	case *map[string]interface{}:
		return *v
	}

	panic(("vjson.deref got unknown type %T"))
}

// unreader is implemented by bytes.Reader and bytes.Buffer
type unreader interface {
	Read(p []byte) (n int, err error)
	// ReadBytes(delim byte) (line []byte, err error)
	ReadByte() (byte, error)
	UnreadByte() error
}

// NOTE: for writing io.Writer works, but for reading io.Reader does NOT work because
// not all JSON data types have a termination character (e.g. you cannot tell when you've
// reached the end of a number without reading past it).  One solution could be to define
// an interface with the methods we need from bytes.Reader, minimally Read() and UnreadByte()

// type MarshalerTo interface {
// 	MarshalJSONTo(w io.Writer) error
// }

// Marshaler is the interface implemented by types that can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// type UnmarshalerFrom interface {
// 	UnmarshalJSONFrom(r io.Reader) error
// }

// Unmarshaler is the interface implemented by types that can unmarshal a JSON description of themselves.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("vjson.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

var _ Marshaler = (*RawMessage)(nil)
var _ Unmarshaler = (*RawMessage)(nil)
