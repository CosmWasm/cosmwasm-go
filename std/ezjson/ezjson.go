package ezjson

import (
	"github.com/cosmwasm/cosmwasm-go/std/ezjson/ext/vjson"
)

func Marshal(v interface{}) ([]byte, error) {
	return vjson.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return vjson.Unmarshal(data, v)
}

func MarshalEx(v interface{}) ([]byte, error) {
	return vjson.MarshalEx(v)
}

func UnmarshalEx(data []byte, vin interface{}) error {
	return vjson.UnmarshalEx(data, vin)
}
