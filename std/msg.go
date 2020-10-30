package std

import (
	"fmt"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// Set this as EmptyStruct{Seen: true} so it will serialize, otherwise it is missing
type EmptyStruct struct {
	Seen bool `json:"seen,omitempty"`
}

var _ ezjson.EzJsonUnmarshaller = EmptyStruct{}

func (e EmptyStruct) UnmarshalEzJson(opts []ezjson.BaseOpt) (interface{}, error) {
	fmt.Printf("Opts: %d\n", len(opts))
	for _, opt := range opts {
		fmt.Printf("Opt: %#v\n", opt)
	}
	return EmptyStruct{Seen: true}, nil
}

// always serialize as an empty struct
func (e EmptyStruct) MarshalJSON() ([]byte, error) {
	return []byte(`{}`), nil
}

// we store seen
func (e *EmptyStruct) UnmarshalJSON(data []byte) error {
	e.Seen = true
	return nil
}
