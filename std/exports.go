package std

import (
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
	"unsafe"
)

func Build_ErrResponse(msg string) string {
	return `{"generic_err":{"msg":"` + msg + `","backtrace":null}}`
}

func StdErrResult(msg string) unsafe.Pointer {
	return Package_message([]byte(`{"Err":` + Build_ErrResponse(msg) + `}`))
}

// =========== Extern --> context =======
type Extern struct {
	EStorage Storage
	EApi     Api
	EQuerier Querier
}

func (deps Extern) change_querier(transform func(Querier) Querier) Extern {
	return Extern{
		EStorage: deps.EStorage,
		EApi:     deps.EApi,
		EQuerier: transform(deps.EQuerier),
	}
}

func make_dependencies() Extern {
	return Extern{
		EStorage: ExternalStorage{},
		EApi:     ExternalApi{},
		EQuerier: ExternalQuerier{},
	}
}

// ========== init ==============
func DoInit(initFn func(deps *Extern, _env Env, msg []byte) (*InitResultOk, *CosmosResponseError), envPtr, msgPtr uint32) unsafe.Pointer {
	var data []byte
	var err error
	envData := TranslateToSlice(uintptr(envPtr))
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, ers := _do_init(initFn, envData, msgData)
	if ok != nil {
		data, err = ezjson.Marshal(*ok)
	} else {
		b, e := ezjson.Marshal(ers)
		if e != nil {
			return StdErrResult("Marshal error failed : " + e.Error())
		}
		return Package_message(b)
	}
	if err != nil {
		return StdErrResult("Failed to marshal init response to []byte: " + err.Error())
	}

	return Package_message(data)
}

func _do_init(initFn func(deps *Extern, _env Env, msg []byte) (*InitResultOk, *CosmosResponseError), envData, msgData []byte) (*InitResultOk, *CosmosResponseError) {
	env := Env{}
	err := ezjson.Unmarshal(envData, &env)
	if err != nil {
		return nil, GenerateError(GenericError, "Testing generic error result", "")
	}
	deps := make_dependencies()
	return initFn(&deps, env, msgData)
}

// ========= handler ============
func DoHandler(handlerFn func(deps *Extern, _env Env, msg []byte) (*HandleResultOk, *CosmosResponseError), envPtr, msgPtr uint32) unsafe.Pointer {
	envData := TranslateToSlice(uintptr(envPtr))
	msgData := Translate_range_custom(uintptr(msgPtr))
	var data []byte
	var err error
	Ok, Err := _do_handler(handlerFn, envData, msgData)
	if Ok != nil {
		data, err = ezjson.Marshal(*Ok)
		if err != nil {
			return StdErrResult("Failed to marshal handle Ok response to []byte: " + err.Error())
		}
	} else if Err != nil {
		data, err = ezjson.Marshal(*Err)
		if err != nil {
			return StdErrResult("Failed to marshal handle Err response to []byte: " + err.Error())
		}
	} else {
		return StdErrResult("Both Ok and Err are nil")
	}
	return Package_message(data)
}

func _do_handler(handlerFn func(deps *Extern, _env Env, msg []byte) (*HandleResultOk, *CosmosResponseError), envData, msgData []byte) (*HandleResultOk, *CosmosResponseError) {
	var env Env
	if err := ezjson.Unmarshal(envData, &env); err != nil {
		return nil, GenerateError(GenericError, "Testing generic error result", "")
	}

	deps := make_dependencies()
	return handlerFn(&deps, env, msgData)
}

// =========== query ===================
func DoQuery(queryFn func(deps *Extern, msg []byte) (*QueryResponseOk, *CosmosResponseError), msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	var data []byte
	var err error
	Ok, Err := _do_query(queryFn, msgData)
	if Ok != nil {
		data, err = ezjson.Marshal(*Ok)
		if err != nil {
			return StdErrResult("Failed to marshal handle Ok response to []byte: " + err.Error())
		}
	} else if Err != nil {
		data, err = ezjson.Marshal(*Err)
		if err != nil {
			return StdErrResult("Failed to marshal handle Err response to []byte: " + err.Error())
		}
	} else {
		return StdErrResult("Both Ok and Err are nil")
	}
	return Package_message(data)
}

func _do_query(handlerFn func(deps *Extern, msg []byte) (*QueryResponseOk, *CosmosResponseError), msgData []byte) (*QueryResponseOk, *CosmosResponseError) {
	deps := make_dependencies()
	return handlerFn(&deps, msgData)
}

//export cosmwasm_vm_version_2
func cosmwasm_vm_version_2() {}

//export allocate
func allocate(size uint32) unsafe.Pointer {
	ptr, _ := Build_region(size, 0)
	return ptr
}

//export deallocate
func deallocate(pointer unsafe.Pointer) {
	Deallocate(pointer)
}
