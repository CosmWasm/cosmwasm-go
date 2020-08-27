package ezjson

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshal(t *testing.T) {

	type TestSt struct {
		Point int
		Key   string
	}
	trans_msg := TestSt{}
	json_msg := "{\"Point\":10086,\"Key\":\"9990293023\"}"
	err := Unmarshal([]byte(json_msg), &trans_msg)
	require.Nil(t, err)
	v := TestSt{
		Point: 100,
		Key:   "123456",
	}
	recv := TestSt{}
	b, e := MarshalEx(v)
	if e == nil {
		str := string(b)
		require.NotNil(t, str)
		e = UnmarshalEx(b, &recv)
		require.Nil(t, e)
	}
}

func TestMarshalA(t *testing.T) {

	type Coin struct {
		Denom  string `json:"denom"`  // type, eg. "ATOM"
		Amount string `json:"amount"` // string encoing of decimal value, eg. "12.3456"
	}

	type TestC struct {
		Point3 int
		Key    string
		Data   byte
		Data2  []byte
	}
	type TestSt struct {
		Point int
		Key   string
		Test3 TestC
	}
	type TestB struct {
		Point2 int
		Key    string
		Test1  TestSt
		Cn     *Coin
	}

	v := TestB{
		Point2: 100,
		Key:    "123456",
		Test1: TestSt{
			Point: 100,
			Key:   "0020202",
			Test3: TestC{
				Point3: 333,
				Key:    "003030303",
				Data:   '1',
				Data2:  []byte("1234567890"),
			},
		},
		Cn: &Coin{
			Denom:  "OOO",
			Amount: "10000202",
		},
	}

	var c []Coin
	c = append(c, Coin{
		Denom:  "OKB",
		Amount: "1000000",
	})
	c = append(c, Coin{
		Denom:  "BTC",
		Amount: "1000000",
	})
	c = append(c, Coin{
		Denom:  "ETH",
		Amount: "1000000",
	})

	var cp *[]Coin
	cp = &c

	b, e := MarshalA(v)
	require.NotNil(t, b)
	require.Nil(t, e)
	fmt.Println(string(b))

	bc, e := MarshalA(cp)
	require.NotNil(t, bc)
	require.Nil(t, e)
	fmt.Println(string(bc))

}
