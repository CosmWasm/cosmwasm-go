package vjson

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return unmarshal(data, v)
}
