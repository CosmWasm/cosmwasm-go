package src

import "C"
import (
	std "../std"
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
	return std.Package_message([]byte("\"Ok\":{\"data\":null,\"log\":[],\"messages\":[]}"))
}

//export handle
func handle(env_ptr uint32, msg_ptr uint32) uint32  {

	return 0
}

//export query
func query(env_ptr uint32, msg_ptr uint32) uint32  {
	return 0
}

func DoNothing(){
	return
}