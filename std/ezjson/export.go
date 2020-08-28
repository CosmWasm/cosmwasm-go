package ezjson

var logging func([]byte) int

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
