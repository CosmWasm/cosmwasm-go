package ezjson

var logging func([]byte) int

func Log(msg string) int {
	if logging != nil {
		logging([]byte(msg))
	}
	return 0
}

func MarshalA(in interface{}) ([]byte, error) {
	return encode2json(in)
}

func SetDisplay(f func([]byte) int) {
	logging = f
}
