package vjson

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return unmarshal(data, v)
}

func MarshalEx(v interface{}) ([]byte, error) {
	return marshalEx(v)
}

func UnmarshalEx(data []byte, v interface{}) error {
	return unmarshalStruct(data, v)
}