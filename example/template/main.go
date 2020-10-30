package main

import (
	"unsafe"

	"github.com/cosmwasm/cosmwasm-go/example/hackatom/src"
	"github.com/cosmwasm/cosmwasm-go/std"
)

func main() {}

//export init
func initialize(env_ptr, info_ptr, msg_ptr uint32) unsafe.Pointer {
	return std.DoInit(src.Init, env_ptr, info_ptr, msg_ptr)
}

//export handle
func handle(env_ptr, info_ptr, msg_ptr uint32) unsafe.Pointer {
	return std.DoHandler(src.Handle, env_ptr, info_ptr, msg_ptr)
}

//export query
func query(env_ptr, msg_ptr uint32) unsafe.Pointer {
	return std.DoQuery(src.Query, env_ptr, msg_ptr)
}
