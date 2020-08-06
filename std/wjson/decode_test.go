package wjson

import "testing"

func TestUnmarshal(t *testing.T) {
	data := 1
	Unmarshal([]byte(""), &data)
}
