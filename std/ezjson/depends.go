package ezjson

var tinygo_typeof = []string{
	"Invalid", "Bool", "Int", "Int8", "Int16", "Int32", "Int64",
	"Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr",
	"Float32", "Float64", "Complex64", "Complex128", "String",
	"UnsafePointer", "Chan", "Interface", "Ptr", "Slice", "Array",
	"Func", "Map", "Struct",
}

func bytesToUint64(b []byte) uint64 {
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func bytesToInt64(buf []byte) int64 {
	return int64(bytesToUint64(buf))
}

func bytesToString(buf []byte) string {
	return string(buf)
}

func bytesToBoolean(buf []byte) bool {
	return true
}

//support flag:
//json:"xxx"
//json:"omitempty"
//json:"xxx,omitempty"
//return name and isOmitEmpty
func getTag(orgTags string) (string, bool) {
	if len(orgTags) < 6 {
		return orgTags, false
	}

	prefix := orgTags[0:6]
	if prefix != "json:\"" {
		return orgTags, false
	}
	begin := 6
	name := ""
	omit := false
	for i, c := range orgTags[6:] {
		if c == 34 { //"
			str := orgTags[begin : i+6]
			if str == "omitempty" {
				omit = true
			} else if len(name) <= 0 {
				name = str
			}
			break
		}
		if c == 44 { //,
			if orgTags[begin:begin+i] == "omitempty" {
				omit = true
			} else {
				name = orgTags[begin : begin+i]
			}
			begin += i + 1 //skip `,`
		}
	}
	return name, omit
}
