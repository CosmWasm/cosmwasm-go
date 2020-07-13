// +build cosmwasm

package src

import "C"
import (
	"github.com/cosmwasm/cosmwasm-go/poc/std"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezdb"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
	"unsafe"
)


//export allocate
func allocate(size uint32) unsafe.Pointer{
	ptr,_ := std.Build_region(size,0)
	return ptr
}

//export deallocate
func deallocate(pointer unsafe.Pointer){
	std.Deallocate(pointer)
}

//export init
func initialize(env_ptr uint32, msg_ptr uint32) unsafe.Pointer {
	var initMsg InitMsg
	msg_data := std.Translate_range_custom(uintptr(msg_ptr))
	str := string(msg_data)
	e := ezjson.UnmarshalEx([]byte(str),&initMsg)
	if e == nil {
		ok,err := go_init(initMsg)
		if ok != nil {
			ezdb.WriteStorage([]byte("inited"),[]byte("true"))
			return std.Package_message([]byte(ok.WrapMessage(ok.Ok)))
		}else {
			return std.Package_message([]byte(err.WrapMessage(err.Err)))
		}
	}else {
		return std.Package_message([]byte(newErrResponse("").WrapMessage(std.Build_ErrResponse("Json analyze failed: " + e.Error()))))
	}
}

//export handle
func handle(env_ptr uint32, msg_ptr uint32) unsafe.Pointer  {
	var handleMsg HandleMsg
	_ ,e:= ezdb.ReadStorage([]byte("inited"))
	if e != nil{
		return std.Package_message([]byte(newErrResponse("").WrapMessage(std.Build_ErrResponse("Uninited contract, need init first"))))
	}
	msg_data := std.Translate_range_custom(uintptr(msg_ptr))
	str := string(msg_data)
	e = ezjson.UnmarshalEx([]byte(str),&handleMsg)
	if e == nil{
		ok,err := go_handle(handleMsg)
		if ok != nil {
			return std.Package_message([]byte(ok.WrapMessage(ok.Ok)))
		}else {
			return std.Package_message([]byte(err.WrapMessage(err.Err)))
		}
	}else {
		return std.Package_message([]byte(newErrResponse("").WrapMessage(std.Build_ErrResponse("Json analyze failed: " + e.Error()))))
	}
}

//export query
func query(msg_ptr uint32) unsafe.Pointer  {
	var queryMsg QueryMsg
	_ ,e:= ezdb.ReadStorage([]byte("inited"))
	if e != nil{
		return std.Package_message([]byte(newErrResponse("").WrapMessage(std.Build_ErrResponse("Uninited contract, need init first"))))
	}
	msg_data := std.Translate_range_custom(uintptr(msg_ptr))
	str := string(msg_data)
	e = ezjson.UnmarshalEx([]byte(str),&queryMsg)
	if e == nil{
		ok,err := go_query(queryMsg)
		if ok != nil {
			return std.Package_message([]byte(ok.WrapMessage(ok.Ok)))
		}else {
			return std.Package_message([]byte(err.WrapMessage(err.Err)))
		}
	}else {
		return std.Package_message([]byte(newErrResponse("").WrapMessage(std.Build_ErrResponse("Json analyze failed: " + e.Error()))))
	}
}

func DoNothing(){
	return
}