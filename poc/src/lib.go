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
	if str == "initialize"{
		ezdb.WriteStorage([]byte("inited"),[]byte("1"))
		return std.Package_message([]byte("{\"Ok\":{\"messages\":[],\"log\":[],\"data\":null}}"))
	}
	return std.Package_message([]byte("{\"Err\":{\"parse_err\":{\"target\":\"InitMsg\",\"msg\":\"you just input a error command, try input initialize\"}}}"))

	//todo need confirm that is reflect package have full support on wasm?
	//https://github.com/tinygo-org/tinygo/issues/1207 waiting response from tinyGo
	e := ezjson.Unmarshal(msg_data,&initMsg)
	if e == nil {
		ok,err := go_init(initMsg)
		if ok != nil {
			b,e := ezjson.Marshal(ok)
			if e == nil {
				return std.Package_message(b)
			}else {
				return std.Package_message([]byte("ezjson.Marshal(ok)" + e.Error()))
			}
		}else {
			b,e := ezjson.Marshal(err)
			if e == nil {
				return std.Package_message(b)
			}else {
				return std.Package_message([]byte("ezjson.Marshal(err)" + e.Error()))
			}
		}
	}else {
		return std.Package_message([]byte("first + " + e.Error()))
	}
}

//export handle
func handle(env_ptr uint32, msg_ptr uint32) unsafe.Pointer  {
	_ ,e:= ezdb.ReadStorage([]byte("inited"))
	if e != nil{
		return std.Package_message([]byte("{\"Err\":{\"messages\":[Uninit contract],\"log\":[],\"data\":null}}"))
	}
	msg_data := std.Translate_range_custom(uintptr(msg_ptr))
	str := string(msg_data)
	ezdb.WriteStorage([]byte("key"),[]byte(str))
	return std.Package_message([]byte("{\"Ok\":{\"messages\":[],\"log\":[],\"data\":null}}"))
}

//export query
func query(msg_ptr uint32) unsafe.Pointer  {
	_ ,e:= ezdb.ReadStorage([]byte("inited"))
	if e != nil{
		return std.Package_message([]byte(std.Build_query_response(std.FakeQueryJson("contract is uninit~ "))))
	}
	v ,e := ezdb.ReadStorage([]byte("key"))
	if e == nil {
		str := string(v[:])
		return std.Package_message([]byte(std.Build_query_response(std.FakeQueryJson(str))))
	}
	return std.Package_message([]byte(std.Build_query_response(e.Error())))
}

func DoNothing(){
	return
}