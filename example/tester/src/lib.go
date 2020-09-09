package src

import "C"
import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"unsafe"
)

//export init
func initialize(env_ptr uint32, msg_ptr uint32) unsafe.Pointer {
	return std.DoInit(Init, env_ptr, msg_ptr)
}

//export handle
func handle(env_ptr uint32, msg_ptr uint32) unsafe.Pointer {
	return std.DoHandler(Invoke, env_ptr, msg_ptr)
}

//export query
func query(msg_ptr uint32) unsafe.Pointer {
	return std.DoQuery(Query, msg_ptr)
}
