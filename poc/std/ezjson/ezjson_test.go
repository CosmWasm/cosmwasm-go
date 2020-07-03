package ezjson

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshal(t *testing.T) {
	type TestSt struct {
		Point int
		Key string
	}
	v := TestSt{
		Point: 100,
		Key: "123456",
	}
	recv := TestSt{}
	b , e := Marshal(v)
	if e == nil {
		str := string(b)
		require.NotNil(t,str)
		e = Unmarshal(b,&recv)
		require.Nil(t,e)
	}

}
