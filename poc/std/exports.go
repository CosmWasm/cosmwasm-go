package std

import (
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
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
func DoInit(initFn func(deps *Extern, _env Env, msg []byte) CosmosResponse, envPtr, msgPtr uint32) unsafe.Pointer {
	envData := Translate_range_custom(uintptr(envPtr))
	msgData := Translate_range_custom(uintptr(msgPtr))

	result := _do_init(initFn, envData, msgData)

	// TODO: you should use json.Marshal, but it does't work now!
	//data, err := ezjson.Marshal(result)
	//if err != nil {
	//	return StdErrResult("Failed to marshal init response to []byte: " + err.Error())
	//}

	data, err := result.MarshalJSON()
	if err != nil {
		return StdErrResult("Failed to marshal init response to []byte: " + err.Error())
	}

	return Package_message(data)
}

func _do_init(initFn func(deps *Extern, _env Env, msg []byte) CosmosResponse, envData, msgData []byte) CosmosResponse {
	var env Env
	err := ezjson.Unmarshal(envData, &env)
	if err != nil {
		return CosmosResponse{
			Ok:  nil,
			Err: ToStdError(GenericErr{err.Error()}),
		}
	}

	deps := make_dependencies()
	return initFn(&deps, env, msgData)
}

// ========= handler ============
func DoHandler(handlerFn func(deps *Extern, _env Env, msg []byte) CosmosResponse, envPtr, msgPtr uint32) unsafe.Pointer {
	envData := Translate_range_custom(uintptr(envPtr))
	msgData := Translate_range_custom(uintptr(msgPtr))

	result := _do_handler(handlerFn, envData, msgData)
	data, err := ezjson.Marshal(result)
	if err != nil {
		return StdErrResult("Failed to marshal handle response to []byte: " + err.Error())
	}

	return Package_message(data)
}

func _do_handler(handlerFn func(deps *Extern, _env Env, msg []byte) CosmosResponse, envData, msgData []byte) CosmosResponse {
	var env Env
	if err := ezjson.Unmarshal(envData, &env); err != nil {
		return CosmosResponse{
			Ok:  nil,
			Err: ToStdError(GenericErr{err.Error()}),
		}
	}

	deps := make_dependencies()
	return handlerFn(&deps, env, msgData)
}

// =========== query ===================
func DoQuery(queryFn func(deps *Extern, msg []byte) CosmosResponse, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))

	result := _do_query(queryFn, msgData)
	data, err := ezjson.Marshal(result)
	if err != nil {
		return StdErrResult("Failed to marshal query response to []byte: " + err.Error())
	}

	return Package_message(data)
}

func _do_query(handlerFn func(deps *Extern, msg []byte) CosmosResponse, msgData []byte) CosmosResponse {
	deps := make_dependencies()
	return handlerFn(&deps, msgData)
}
