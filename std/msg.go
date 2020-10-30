package std

import (
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// Set this as EmptyStruct{Seen: true} so it will serialize, otherwise it is missing
type EmptyStruct struct {
	Seen bool `json:"do_not_set_this_field,omitempty"`
}

var _ ezjson.EzJsonUnmarshaller = EmptyStruct{}

func (e EmptyStruct) UnmarshalEzJson(opts []ezjson.BaseOpt) (interface{}, error) {
	// Odd but true - if struct was seen with no data, len(opts) == 0
	// If it was not seen, len(opts) == 1, tagged with the default name ("do_not_set_this_field")
	// If it was a struct with some data, len(opts) >= 1, tagged with the actual data present eg ({"a": 1} => "a")
	unseen := len(opts) == 1 && opts[0].Tag() == "do_not_set_this_field"
	return EmptyStruct{Seen: !unseen}, nil
}

func (e EmptyStruct) WasSet() bool {
	return e.Seen
}

//// always serialize as an empty struct
//func (e EmptyStruct) MarshalJSON() ([]byte, error) {
//	return []byte(`{}`), nil
//}
