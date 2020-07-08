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
	trans_msg := TestSt{}
	json_msg := "{\"Point\":10086,\"Key\":\"9990293023\"}"
	err := Unmarshal([]byte(json_msg),&trans_msg)
	require.Nil(t,err)
	v := TestSt{
		Point: 100,
		Key: "123456",
	}
	recv := TestSt{}
	b , e := MarshalEx(v)
	if e == nil {
		str := string(b)
		require.NotNil(t,str)
		e = UnmarshalEx(b,&recv)
		require.Nil(t,e)
	}
}

