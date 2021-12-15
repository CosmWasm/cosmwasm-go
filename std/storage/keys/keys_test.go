package keys

import (
	"bytes"
	"reflect"
	"testing"
)

type comparison int

const (
	comparisonSmaller comparison = -1
	comparisonEqual   comparison = 0
	comparisonBigger  comparison = 1
)

func TestInt32PrimaryKey(t *testing.T) {
	type test struct {
		v1, v2  int32
		compare comparison
	}

	tests := map[string]test{
		"negative < positive": {
			v1:      -100,
			v2:      100,
			compare: comparisonSmaller,
		},
		"positive equality": {
			v1:      100,
			v2:      100,
			compare: comparisonEqual,
		},

		"negative equality": {
			v1:      -100,
			v2:      -100,
			compare: comparisonEqual,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {

			compare(t, Int32PrimaryKey, tc.v1, tc.v2, tc.compare)
		})
	}
}

// compare takes a byte producer function such as Int32PrimaryKey, StringPrimaryKey
// two valid values, and compares the values. Fails if comparison does not match expectations.
func compare(t *testing.T, byteProducer, v1, v2 interface{}, comp comparison) {
	human := map[int]string{
		-1: "smaller",
		0:  "equal",
		1:  "bigger",
	}

	f := reflect.ValueOf(byteProducer)
	v1v := reflect.ValueOf(v1)
	v2v := reflect.ValueOf(v2)

	v1Bytes := f.Call([]reflect.Value{v1v})[0].Bytes()
	v2Bytes := f.Call([]reflect.Value{v2v})[0].Bytes()

	got := bytes.Compare(v1Bytes, v2Bytes)
	if got == int(comp) {
		return
	}

	expectedHuman := human[int(comp)]
	gotHuman := human[got]

	t.Fatalf("test expected %v to be %s than %v, got: %s", v1, expectedHuman, v2, gotHuman)
}
