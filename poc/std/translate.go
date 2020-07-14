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
	"encoding/base64"
	"unsafe"
)

func Build_region(size uint32, len uint32) (unsafe.Pointer, *MemRegion) {
	ptr := C.malloc(C.ulong(size) + C.ulong(REGION_HEAD_SIZE))
	var region = new(MemRegion)
	region.Offset = uint32(uintptr(ptr)) + REGION_HEAD_SIZE
	region.Capacity = uint32(size)
	region.Length = len
	C.memcpy(ptr, unsafe.Pointer(region), C.ulong(REGION_HEAD_SIZE))
	return ptr, region
}

func Translate_range_custom(ptr uintptr) []byte {
	if ptr == 0 {
		return nil
	}
	var mm []byte
	region := (*MemRegion)(unsafe.Pointer(ptr))
	header := (*SliceHeader_tinyGo)(unsafe.Pointer(&mm))

	header.Len = uintptr(region.Length)
	header.Cap = uintptr(region.Capacity)
	header.Data = uintptr(region.Offset)
	return mm
}

func TranslateToSlice(ptr uintptr) []byte {
	if ptr == 0 {
		return nil
	}
	region := (*MemRegion)(unsafe.Pointer(ptr))
	header := SliceHeader_tinyGo{
		Data: ptr + 12,
		Len:  uintptr(region.Length),
		Cap:  uintptr(region.Capacity),
	}
	b := *(*[]byte)(unsafe.Pointer(&header))
	return b
}

func TranslateToRegion(b []byte, ptr uintptr) uintptr {
	if b == nil || ptr == 0 {
		return 0
	}
	header := (*SliceHeader_tinyGo)(unsafe.Pointer(&b))
	region := (*MemRegion)(unsafe.Pointer(ptr))

	region.Length = uint32(header.Len)
	region.Capacity = uint32(header.Cap)
	region.Offset = uint32(header.Data)
	return ptr
}

func Deallocate(pointer unsafe.Pointer) {
	C.free(pointer)
}

func Package_message(msg []byte) unsafe.Pointer {
	size := len(msg)
	ptr, _ := Build_region(uint32(size), uint32(size))
	result := uintptr(ptr) + uintptr(REGION_HEAD_SIZE)
	for _, m := range msg {
		*(*byte)(unsafe.Pointer(result)) = byte(m)
		result += 1
	}
	return ptr
}

func Build_QueryResponse(msg string) string {
	encoding := base64.StdEncoding.EncodeToString([]byte(`{"QueryResult":"` + msg + `"}`))
	return `"` + encoding + `"`
}

func Build_OkResponse(msg string) string {
	return `{"messages":[],"log":[{"key":"result","value":"` + msg + `"}],"data":null}`
}

func Build_ErrResponse(msg string) string {
	return `{"generic_err":{"msg":"` + msg + `","backtrace":null}}`
}
