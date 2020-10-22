package main

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/example/tester/src"
	"unsafe"
)

func main() {}

//export init
func initialize(env_ptr uint32, msg_ptr uint32) unsafe.Pointer {
	return std.DoInit(src.Init, env_ptr, msg_ptr)
}

//export handle
func handle(env_ptr uint32, msg_ptr uint32) unsafe.Pointer {
	return std.DoHandler(src.Invoke, env_ptr, msg_ptr)
}

//export query
func query(msg_ptr uint32) unsafe.Pointer {
	return std.DoQuery(src.Query, msg_ptr)
}
