package wjson

import (
	"bytes"
	"encoding/base64"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var hex = "0123456789abcdef"

type encoder struct {
	bytes.Buffer
	scratch [64]byte
}

func newEncoder() *encoder {
	return new(encoder)
}

func (e *encoder) marshal(v interface{}, opts encOpts) error {
	return e.reflectValue(reflect.ValueOf(v), opts)
}

func (e *encoder) reflectValue(v reflect.Value, opts encOpts) error {
	return valueEncoder(v)(e, v, opts)
}

func (e *encoder) string(s string, escapeHTML bool) error {
	_ = e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				_, _ = e.WriteString(s[start:i])
			}
			_ = e.WriteByte('\\')
			switch b {
			case '\\', '"':
				_ = e.WriteByte(b)
			case '\n':
				_ = e.WriteByte('n')
			case '\r':
				_ = e.WriteByte('r')
			case '\t':
				_ = e.WriteByte('t')
			default:
				_, _ = e.WriteString(`u00`)
				_ = e.WriteByte(hex[b>>4])
				_ = e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				_, _ = e.WriteString(s[start:i])
			}
			_, _ = e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}

		if c == '\u2028' || c == '\u2029' {
			if start < i {
				_, _ = e.WriteString(s[start:i])
			}
			_, _ = e.WriteString(`\u202`)
			_ = e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		_, _ = e.WriteString(s[start:])
	}
	_ = e.WriteByte('"')
	return nil
}

type encOpts struct {
	quoted     bool
	escapeHTML bool
}

type encoderFunc func(e *encoder, v reflect.Value, opts encOpts) error

// map[reflect.Type]encoderFunc
var encoderCache sync.Map

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}

	return typeEncoder(v.Type())
}

func invalidValueEncoder(e *encoder, _ reflect.Value, _ encOpts) error {
	_, err := e.WriteString("null")
	return err
}

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	f := newTypeEncoder(t)
	encoderCache.Store(t, f)
	return f
}

func newTypeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	//case reflect.Float32:
	//	return float32Encoder
	//case reflect.Float64:
	//	return float64Encoder
	case reflect.String:
		return stringEncoder
	//case reflect.Interface:
	//	return interfaceEncoder
	//TODO
	case reflect.Struct:
		return newStructEncoder(t)
	//case reflect.Map:
	//	return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	//case reflect.Ptr:
	//	return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func boolEncoder(e *encoder, v reflect.Value, opts encOpts) (err error) {
	if opts.quoted {
		err = e.WriteByte('"')
	}
	if v.Bool() {
		_, err = e.WriteString("true")
	} else {
		_, err = e.WriteString("false")
	}
	if opts.quoted {
		err = e.WriteByte('"')
	}
	return
}

func intEncoder(e *encoder, v reflect.Value, opts encOpts) (err error) {
	b := strconv.AppendInt(e.scratch[:0], v.Int(), 10)
	if opts.quoted {
		err = e.WriteByte('"')
	}
	_, err = e.Write(b)
	if opts.quoted {
		err = e.WriteByte('"')
	}
	return
}

func uintEncoder(e *encoder, v reflect.Value, opts encOpts) (err error) {
	b := strconv.AppendUint(e.scratch[:0], v.Uint(), 10)
	if opts.quoted {
		err = e.WriteByte('"')
	}
	_, err = e.Write(b)
	if opts.quoted {
		err = e.WriteByte('"')
	}
	return
}

func stringEncoder(e *encoder, v reflect.Value, opts encOpts) (err error) {
	if opts.quoted {
		sb, err := Marshal(v.String())
		if err != nil {
			return errors.New("JSON marshal error: string type")
		}
		err = e.string(string(sb), opts.escapeHTML)
	} else {
		err = e.string(v.String(), opts.escapeHTML)
	}
	return
}

func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{fields: cachedTypeFields(t)}
	return se.encode
}

type structEncoder struct {
	fields []field
}

func (se structEncoder) encode(e *encoder, v reflect.Value, opts encOpts) (err error) {
	next := byte('{')
FieldLoop:
	for i := range se.fields {
		f := &se.fields[i]
		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					continue FieldLoop
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}

		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		_ = e.WriteByte(next)
		next = ','
		if opts.escapeHTML {
			_, _ = e.WriteString(f.nameEscHTML)
		} else {
			_, _ = e.WriteString(f.nameNonEsc)
		}
		opts.quoted = f.quoted
		if err = f.encoder(e, fv, opts); err != nil {
			return err
		}
	}
	if next == '{' {
		_, err = e.WriteString("{}")
	} else {
		err = e.WriteByte('}')
	}
	return err
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

type field struct {
	name        string
	nameBytes   []byte
	equalFold   func(s, t []byte) bool
	nameNonEsc  string
	nameEscHTML string
	tag         bool
	index       []int
	typ         reflect.Type
	omitEmpty   bool
	quoted      bool
	encoder     encoderFunc
}

// map[reflect.Type][]field
var fieldCache sync.Map

func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]field)
}

type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

func typeFields(t reflect.Type) []field {
	current := []field{}
	next := []field{{typ: t}}
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}
	visited := map[reflect.Type]bool{}
	var (
		fields     []field
		nameEscBuf bytes.Buffer
	)

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}
					if isUnexported && t.Kind() != reflect.Struct {
						continue
					}
				} else if isUnexported {
					continue
				}

				// only support tag like `json:"xxx"`
				tag := getTag(string(sf.Tag), "json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if len(sf.Name) == 0 && ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}

				quoted := false
				if opts.Contains("string") {
					switch ft.Kind() {
					case reflect.Bool,
						reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
						reflect.Float32, reflect.Float64,
						reflect.String:
						quoted = true
					}
				}

				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := field{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),
						quoted:    quoted,
					}
					field.nameBytes = []byte(field.name)
					field.equalFold = foldFunc(field.nameBytes)
					nameEscBuf.Reset()
					nameEscBuf.WriteString(`"`)
					HTMLEscape(&nameEscBuf, field.nameBytes)
					nameEscBuf.WriteString(`":`)
					field.nameEscHTML = nameEscBuf.String()
					field.nameNonEsc = `"` + field.name + `":`
					fields = append(fields, field)
					if count[f.typ] > 1 {
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: sf.Name, index: index, typ: ft})
				}
			}
		}
	}

	fieldsLen := len(fields)
	for i := 0; i < fieldsLen-1; i++ {
		for j := i + 1; j < fieldsLen; j++ {
			if fields[i].name > fields[j].name {
				fields[i], fields[j] = fields[j], fields[i]
			} else if len(fields[i].index) > len(fields[j].index) {
				fields[i], fields[j] = fields[j], fields[i]
			} else if fields[i].tag != fields[j].tag {
				if fields[j].tag {
					fields[i], fields[j] = fields[j], fields[i]
				}
			}

			if !byIndex(fields).Less(i, j) {
				fields[i], fields[j] = fields[j], fields[i]
			}
		}
	}

	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 {
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	for i := range fields {
		f := &fields[i]
		f.encoder = typeEncoder(typeByIndex(t, f.index))
	}
	return fields
}

func getTag(structTag, key string) string {
	for structTag != "" {
		i := 0
		for i < len(structTag) && structTag[i] == ' ' {
			i++
		}
		structTag = structTag[i:]
		if structTag == "" {
			break
		}

		i = 0
		for i < len(structTag) && structTag[i] > ' ' && structTag[i] != ':' && structTag[i] != '"' && structTag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(structTag) || structTag[i] != ':' || structTag[i+1] != '"' {
			break
		}
		name := string(structTag[:i])
		structTag = structTag[i+1:]

		i = 1
		for i < len(structTag) && structTag[i] != '"' {
			if structTag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(structTag) {
			break
		}
		qvalue := string(structTag[:i+1])
		structTag = structTag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value
		}
	}
	return ""
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	if strings.ContainsAny(s, "!#$%&()*+-./:<=>?@[]^_{|}~ ") {
		return false
	}

	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func HTMLEscape(dst *bytes.Buffer, src []byte) {
	start := 0
	for i, c := range src {
		if c == '<' || c == '>' || c == '&' {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u00`)
			dst.WriteByte(hex[c>>4])
			dst.WriteByte(hex[c&0xF])
			start = i + 1
		}
		if c == 0xE2 && i+2 < len(src) && src[i+1] == 0x80 && src[i+2]&^1 == 0xA8 {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u202`)
			dst.WriteByte(hex[src[i+2]&0xF])
			start = i + 3
		}
	}
	if start < len(src) {
		dst.Write(src[start:])
	}
}

func dominantField(fields []field) (field, bool) {
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tag == fields[1].tag {
		return field{}, false
	}
	return fields[0], true
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae arrayEncoder) encode(e *encoder, v reflect.Value, opts encOpts) error {
	_ = e.WriteByte('[')
	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			_ = e.WriteByte(',')
		}
		ae.elemEnc(e, v.Index(i), opts)
	}
	return e.WriteByte(']')
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	//if t.Elem().Kind() == reflect.Uint8 {
	//	p := reflect.PtrTo(t.Elem())
	//	//if !p.Implements(marshalerType) && !p.Implements(textMarshalerType) {
	//	//	return encodeByteSlice
	//	//}
	//}
	enc := sliceEncoder{newArrayEncoder(t)}
	return enc.encode
}

type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se sliceEncoder) encode(e *encoder, v reflect.Value, opts encOpts) error {
	if v.IsNil() {
		_, _ = e.WriteString("null")
		return nil
	}
	return se.arrayEnc(e, v, opts)
}

func encodeByteSlice(e *encoder, v reflect.Value, _ encOpts) error {
	if v.IsNil() {
		e.WriteString("null")
		return nil
	}
	s := v.Bytes()
	_ = e.WriteByte('"')
	encodedLen := base64.StdEncoding.EncodedLen(len(s))
	if encodedLen <= len(e.scratch) {
		// If the encoded bytes fit in e.scratch, avoid an extra
		// allocation and use the cheaper Encoding.Encode.
		dst := e.scratch[:encodedLen]
		base64.StdEncoding.Encode(dst, s)
		_, _ = e.Write(dst)
	} else if encodedLen <= 1024 {
		// The encoded bytes are short enough to allocate for, and
		// Encoding.Encode is still cheaper.
		dst := make([]byte, encodedLen)
		base64.StdEncoding.Encode(dst, s)
		_, _ = e.Write(dst)
	} else {
		// The encoded bytes are too long to cheaply allocate, and
		// Encoding.Encode is no longer noticeably cheaper.
		enc := base64.NewEncoder(base64.StdEncoding, e)
		_, _ = enc.Write(s)
		_ = enc.Close()
	}
	return e.WriteByte('"')
}

func unsupportedTypeEncoder(e *encoder, _ reflect.Value, _ encOpts) error {
	return errors.New("JSON marshal error: unsupported type")
}
