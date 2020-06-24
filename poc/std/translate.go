package std
//#include "stdlib.h"
//#include "string.h"
/*
void CopyMessage(char* to,char* from,unsigned int size){
	int i = 0;
	while(i++ < size){
		to[i] = from[i];
	}
}
 */
import "C"
import (
	"unsafe"
)

func Build_region(size uint32,len uint32) (unsafe.Pointer,*MemRegion){
	ptr := C.malloc(C.ulong(size) + C.ulong(REGION_HEAD_SIZE))
	var region = new(MemRegion)
	region.Offset = uint32(uintptr(ptr)) + REGION_HEAD_SIZE
	region.Capacity = uint32(size)
	region.Length = len
	C.memcpy(ptr,unsafe.Pointer(region),C.ulong(REGION_HEAD_SIZE))
	return ptr,region
}

func Deallocate(pointer unsafe.Pointer){
	C.free(pointer)
}

func Package_message(msg []byte) unsafe.Pointer{
	size := len(msg)
	ptr,_ := Build_region(uint32(size),uint32(size))
	result := uintptr(ptr) + uintptr(REGION_HEAD_SIZE)
	for _,m := range msg {
		*(*byte)(unsafe.Pointer(result)) = byte(m)
		result += 1
	}
	return ptr
}
