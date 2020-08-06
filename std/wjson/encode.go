package wjson

func Marshal(v interface{}) ([]byte, error) {
	encoder := newEncoder()
	if err := encoder.marshal(v, encOpts{escapeHTML: true}); err != nil {
		return nil, err
	}
	buffer := append([]byte(nil), encoder.Bytes()...)
	return buffer, nil
}
