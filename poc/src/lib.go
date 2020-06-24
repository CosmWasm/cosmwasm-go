package src

//#include "stdlib.h"
/*
unsigned int allocate(unsigned int size){
	char* addr = malloc(size);
	unsigned int result = (unsigned int)addr;
	return 10086;
}

void deallocate(unsigned int pointer){
	free((char*)pointer);
}

 */
import "C"
import "unsafe"

//export allocate
func allocate(size C.uint) unsafe.Pointer{
	ptr := C.malloc(C.ulong(size))
	return ptr
}

//export deallocate
func deallocate(pointer C.uint){
	C.deallocate(pointer);
}

//export init
func initialize(env_ptr uint32, msg_ptr uint32) uint32 {
	return 0
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