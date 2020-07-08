package ezdb

/*
#include "stdlib.h"
extern int db_read(void* key, void* value);
extern int db_write(void* key, void* value);
extern int db_remove(void* key);

*/
import "C"
import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/poc/std"
	"unsafe"
)

func WriteStorage(key []byte,value []byte) (error){
	ptr := C.malloc(C.ulong(len(key)))
	val := C.malloc(C.ulong(len(value)))
	region_key 	 := std.TranslateToRegion(key,uintptr(ptr))
	region_value := std.TranslateToRegion(value,uintptr(val))
	ret := C.db_write(unsafe.Pointer(region_key),unsafe.Pointer(region_value))
	C.free(ptr)
	C.free(val)
	if ret != 0{
		return errors.New("Some Error happend during Write storage")
	}
	return nil
}

func ReadStorage(key []byte) ([]byte,error){

	ptr := C.malloc(C.ulong(len(key)))
	val_ptr,_ := std.Build_region(1024,0)
	region := std.TranslateToRegion(key,uintptr(ptr))
	ret := C.db_read(unsafe.Pointer(region),unsafe.Pointer(val_ptr))
	C.free(ptr)
	if ret != 0{
		if ret == -1001001{
			return nil,errors.New("key not existed~ ")
		}
		return nil,errors.New("call success but reading failed ")
	}
	b := std.TranslateToSlice(uintptr(val_ptr))
	return b,nil
}