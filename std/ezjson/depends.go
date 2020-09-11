package ezjson

var tinygo_typeof = []string{
	"Invalid", "Bool", "Int", "Int8", "Int16", "Int32", "Int64",
	"Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr",
	"Float32", "Float64", "Complex64", "Complex128", "String",
	"UnsafePointer", "Chan", "Interface", "Ptr", "Slice", "Array",
	"Func", "Map", "Struct",
}

const (
	OmitEmpty  = "omitempty"
	RustOption = "rust_option"
)

//support flag:
//json:"xxx"
//json:"omitempty"
//json:"xxx,omitempty"
//return name and isOmitEmpty,rust_option
func getTag(orgTags string) (string, bool, bool) {
	if len(orgTags) < 6 {
		return orgTags, false, false
	}

	prefix := orgTags[0:6]
	if prefix != "json:\"" {
		return orgTags, false, false
	}
	begin := 6
	name := ""
	omit := false
	rustOption := false
	for i, c := range orgTags[6:] {
		if c == 34 { //"
			str := orgTags[begin : i+6]
			if str == OmitEmpty {
				omit = true
			} else if str == RustOption {
				rustOption = true
			} else if len(name) <= 0 {
				name = str
			}
			break
		}
		if c == 44 { //,
			str := orgTags[begin : i+6]
			if str == OmitEmpty {
				omit = true
			} else if str == RustOption {
				rustOption = true
			} else {
				name = str
			}
			begin = i + 7 //skip `,`
		}
	}
	return name, omit, rustOption
}
