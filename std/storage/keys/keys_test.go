package keys

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type comparison int

func (c comparison) String() string {
	switch c {
	case comparisonSmaller:
		return "smaller"
	case comparisonEqual:
		return "equal"
	case comparisonBigger:
		return "bigger"
	default:
		panic(fmt.Errorf("unknown human representation for %d", c))
	}
}

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

func TestInt64PrimaryKey(t *testing.T) {
	type test struct {
		v1, v2  int64
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

			compare(t, Int64PrimaryKey, tc.v1, tc.v2, tc.compare)
		})
	}
}

func TestInt16PrimaryKey(t *testing.T) {
	type test struct {
		v1, v2  int16
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

			compare(t, Int16PrimaryKey, tc.v1, tc.v2, tc.compare)
		})
	}
}

func TestInt8PrimaryKey(t *testing.T) {
	type test struct {
		v1, v2  int8
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

			compare(t, Int8PrimaryKey, tc.v1, tc.v2, tc.compare)
		})
	}
}

// compare takes a byte producer function such as Int32PrimaryKey, StringPrimaryKey
// two valid values, and compares the values. Fails if comparison does not match expectations.
func compare(t *testing.T, byteProducer, v1, v2 interface{}, comp comparison) {

	f := reflect.ValueOf(byteProducer)
	v1v := reflect.ValueOf(v1)
	v2v := reflect.ValueOf(v2)

	v1Bytes := f.Call([]reflect.Value{v1v})[0].Bytes()
	v2Bytes := f.Call([]reflect.Value{v2v})[0].Bytes()

	got := comparison(bytes.Compare(v1Bytes, v2Bytes))
	if got == comp {
		return
	}

	t.Fatalf("test expected %v to be %s than %v, got: %s", v1, comp, v2, got)
}
