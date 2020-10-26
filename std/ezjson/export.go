package ezjson

import "fmt"

//var logging func([]byte) int = nil
var logging func([]byte) int = func(bz []byte) int {
	fmt.Println(string(bz))
	return len(bz)
}

func Log(msg string) int {
	if logging != nil {
		logging([]byte(msg))
	}
	return 0
}

func Marshal(in interface{}) ([]byte, error) {
	return encode2json(in)
}

func Unmarshal(jsonstr []byte, out interface{}) error {
	return decodeJson(jsonstr, out)
}

func SetDisplay(f func([]byte) int) {
	logging = f
}
