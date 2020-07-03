package ezjson

import (
	"github.com/cosmwasm/cosmwasm-go/poc/ext/vjson"
)

func Marshal(v interface{}) ([]byte, error){
	return vjson.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return vjson.Unmarshal(data,v)
}
