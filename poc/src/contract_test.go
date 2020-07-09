package src

import (
	"fmt"
	"github.com/cosmwasm/cosmwasm-go/poc/std"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAllocate(t *testing.T) {
	result := std.Package_message([]byte("1234567"))
	_ = result
}

func TestInit(t *testing.T){

	msg := InitMsg{
		Verifier:    "AAAAAA",
		Beneficiary: "BBBBBB",
		Legacy:      1100,
	}
	msg_alloc := allocate(1000)
	msg_str,e := ezjson.Marshal(msg)
	require.Nil(t,e)
	b := std.Translate_range_custom(uintptr(msg_alloc))
	//src := std.Translate_to_slice(uintptr(unsafe.Pointer(&msg)),unsafe.Sizeof(msg))
	copy(b,[]byte(msg_str))

	OkResp := OKResponse{}
	result := InitResponse{}
	ret := initialize(0,uint32(uintptr(msg_alloc)))
	if ret == nil{
		bstr ,e:= ezjson.Marshal("")
		require.Nil(t,e)
		fmt.Println(string(bstr))
		e = ezjson.Unmarshal(bstr,&OkResp)
		require.Nil(t,e)
		e = ezjson.Unmarshal([]byte(OkResp.Ok),&result)
		require.Nil(t,e)
	}


}