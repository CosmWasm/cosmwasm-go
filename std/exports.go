// +build cosmwasm

package std

import (
	"strconv"
	"strings"
	"unsafe"

	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

type ContractError struct {
	Err string `json:"error"`
}

func StdErrResult(err error, prefix string) unsafe.Pointer {
	clean := strings.Replace(err.Error(), `"`, `\"`, -1)
	msg := `{"error":"` + clean + `"}`
	//msg := `{"error":"` + prefix + `"}`
	return Package_message([]byte(msg))

	//msg := err.Error()
	//if prefix != "" {
	//	msg = prefix + "- " + msg
	//}
	//e := ContractError{Err: msg}
	//bz, _ := ezjson.Marshal(e)
	//return Package_message(bz)
}

func make_dependencies() Extern {
	return Extern{
		EStorage: ExternalStorage{},
		EApi:     ExternalApi{},
		EQuerier: ExternalQuerier{},
	}
}

// ========== init ==============
func DoInit(initFn func(*Extern, Env, MessageInfo, []byte) (*InitResultOk, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	deps := make_dependencies()

	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := ezjson.Unmarshal(envData, &env)
	if err != nil {
		return StdErrResult(err, "Parse Env")
	}

	info := MessageInfo{
		//SentFunds: []Coin{{}},
		SentFunds: make([]Coin, 5, 10),
	}
	infoData := TranslateToSlice(uintptr(infoPtr))

	err = ezjson.Unmarshal(infoData, &info)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}
	// measure funds length
	deps.EApi.Debug("Len: " + strconv.Itoa(len(info.SentFunds)))

	// remove 0's
	var i = 0
	for info.SentFunds[i].Denom != "" {
		i++
	}
	info.SentFunds = info.SentFunds[:i]
	deps.EApi.Debug("Trimmed Len: " + strconv.Itoa(len(info.SentFunds)))

	ok, err := initFn(&deps, env, info, msgData)
	if ok == nil {
		return StdErrResult(err, "")
	}

	data, err := ezjson.Marshal(*ok)
	if err != nil {
		return StdErrResult(err, "Marshal Response")
	}
	return Package_message(data)
}

// ========= handler ============
func DoHandler(handlerFn func(*Extern, Env, MessageInfo, []byte) (*HandleResultOk, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := ezjson.Unmarshal(envData, &env)
	if err != nil {
		return StdErrResult(err, "Parse Env")
	}
	info := MessageInfo{}
	infoData := TranslateToSlice(uintptr(infoPtr))
	err = ezjson.Unmarshal(infoData, &info)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}

	deps := make_dependencies()
	ok, err := handlerFn(&deps, env, info, msgData)
	if ok == nil {
		return StdErrResult(err, "")
	}

	data, err := ezjson.Marshal(*ok)
	if err != nil {
		return StdErrResult(err, "Marshal Response")
	}
	return Package_message(data)
}

// =========== query ===================
func DoQuery(queryFn func(*Extern, Env, []byte) (*QueryResponseOk, error), envPtr, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := ezjson.Unmarshal(envData, &env)
	if err != nil {
		return StdErrResult(err, "Parse Env")
	}

	deps := make_dependencies()
	ok, err := queryFn(&deps, env, msgData)
	if ok == nil {
		return StdErrResult(err, "")
	}

	data, err := ezjson.Marshal(*ok)
	if err != nil {
		return StdErrResult(err, "Marshal Response")
	}
	return Package_message(data)
}

//export cosmwasm_vm_version_4
func cosmwasm_vm_version_4() {}

//export allocate
func allocate(size uint32) unsafe.Pointer {
	ptr, _ := Build_region(size, 0)
	return ptr
}

//export deallocate
func deallocate(pointer unsafe.Pointer) {
	Deallocate(pointer)
}
