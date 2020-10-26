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
	b, e := Marshal(v)
	if e == nil {
		str := string(b)
		require.NotNil(t, str)
		e = Unmarshal(b, &recv)
		require.Nil(t, e)
	}
}

func TestMarshalA(t *testing.T) {

	type Coin struct {
		Denom  string `denom`  // type, eg. "ATOM"
		Amount string `amount` // string encoing of decimal value, eg. "12.3456"
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
		Cn     Coin
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
		Cn: Coin{
			Denom:  "OOO",
			Amount: "10000202",
		},
	}

	b, e := Marshal(v)
	require.NotNil(t, b)
	require.Nil(t, e)
	fmt.Println(string(b))
}

func TestUnmarshal(t *testing.T) {

	// Coin is a string representation of the sdk.Coin type (more portable than sdk.Int)
	type Coin struct {
		Denom  string `denom`  // type, eg. "ATOM"
		Amount string `amount` // string encoing of decimal value, eg. "12.3456"
	}
	type BlockInfo struct {
		// block height this transaction is executed
		Height uint64 `json:"height"`
		// time in seconds since unix epoch - since cosmwasm 0.3
		Time    uint64 `json:"time"`
		ChainID string `json:"chain_id"`
	}

	type MessageInfo struct {
		// binary encoding of sdk.AccAddress executing the contract
		Sender []byte `json:"sender"`
		// amount of funds send to the contract along with this message
		SentFunds []Coin `json:"sent_funds,omitempty"`
	}

	type ContractInfo struct {
		// binary encoding of sdk.AccAddress of the contract, to be used when sending messages
		Address []byte `json:"address"`
	}
	type Env struct {
		Block    BlockInfo    `json:"block"`
		Message  MessageInfo  `json:"message"`
		Contract ContractInfo `json:"contract,omitempty"`
	}

	str := "{\"block\":{\"height\":12345,\"time\":1571797419,\"chain_id\":\"cosmos-testnet-14002\"},\"message\":{\"sender\":\"original_owner_addr\",\"sent_funds\":[]},\"contract\":{\"address\":\"cosmos2contract\"}}"
	var obj Env
	e := Unmarshal([]byte(str), &obj)
	require.Nil(t, e)
	obj.Contract.Address = nil //set to nil, test omitempty keyword
	b, e := Marshal(obj)
	require.Nil(t, e)
	fmt.Println(string(b))
}

func TestMarshal_Tag(t *testing.T) {
	name, b, _ := getTag("json:\"time\"")
	require.Equal(t, "time", name)
	require.Equal(t, b, false)

	name, b, _ = getTag("json:\"opt,omitempty\"")
	require.Equal(t, "opt", name)
	require.Equal(t, b, true)

	name, b, r := getTag("json:\"omitempty\"")
	require.Equal(t, "", name)
	require.Equal(t, b, true)
	require.Equal(t, r, false)

	name, _, r = getTag("json:\"rust_option\"")
	require.Equal(t, "", name)
	require.Equal(t, r, true)

	name, _, r = getTag(`json:"Ok,omitempty,rust_option"`)
	require.Equal(t, "Ok", name)
	require.Equal(t, r, true)

}

func TestRustOption(t *testing.T) {
	type Coin struct {
		Denom  string `json:"denom,rust_option"` // type, eg. "ATOM"
		Amount string `json:"amount"`            // string encoing of decimal value, eg. "12.3456"
		Name   string
	}
	c := Coin{
		Denom:  "",
		Amount: "1000",
		Name:   "",
	}
	b, e := Marshal(c)
	require.Nil(t, e)
	require.NotNil(t, b)
	fmt.Println(string(b))
}

func TestEnvFailure(t *testing.T) {
	type Coin struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	type MessageInfo struct {
		Sender    string `json:"sender"`
		SentFunds []Coin `json:"sent_funds"`
	}

	sender := "coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w"

	cases := map[string]struct {
		msg      []byte
		expected MessageInfo
	}{
		//"nil coins": {
		//	msg:      []byte(`{"sender":"coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w"}`),
		//	expected: MessageInfo{Sender: sender},
		//},
		//"zero coins": {
		//	msg: []byte(`{"sender":"coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w","sent_funds":[]}`),
		//	// [] decoded as nil
		//	expected: MessageInfo{Sender: sender},
		//},
		"one coin": {
			msg:      []byte(`{"sender":"coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w","sent_funds":[{"denom":"uatom","amount":"1000"}]}`),
			expected: MessageInfo{Sender: sender, SentFunds: []Coin{{Denom: "uatom", Amount: "1000"}}},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var info MessageInfo
			err := Unmarshal(tc.msg, &info)
			require.NoError(t, err)
			require.Equal(t, tc.expected, info)
		})
	}

}
