package std

import (
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
	"unsafe"
)

// =========== Extern --> context =======
type Extern struct {
	EStorage Storage
	EApi Api
	EQuerier Querier
}

func (deps Extern) change_querier(transform func(Querier) Querier) Extern{
	return Extern{
		EStorage: deps.EStorage,
		EApi:     deps.EApi,
		EQuerier: transform(deps.EQuerier),
	}
}

func make_dependencies() Extern {
	return Extern{
		EStorage: ExternalStorage{},
		EApi: ExternalApi{},
		EQuerier: ExternalQuerier{},
	}
}

// ========== init ==============
type InitResponse struct {
	Messages []string `json:"messages"`
	Logs []string `json:"log"`
	Data []byte `json:"data"`
}

func DoInit(initFn func(deps *Extern, _env Env, msg []byte) InitResponse, envPtr, msgPtr uint32) unsafe.Pointer{
	envData := Translate_range_custom(uintptr(envPtr))
	msgData := Translate_range_custom(uintptr(msgPtr))

	result := _do_init(initFn, envData, msgData)

	data, err := ezjson.Marshal(result)
	if err != nil {
		return nil
	}

	return Package_message(data)
}

func _do_init(initFn func(deps *Extern, _env Env, msg []byte) InitResponse, envData, msgData []byte) InitResponse{
	var env Env
	ezjson.Unmarshal(envData, &env)

	deps := make_dependencies()
	return initFn(&deps, env, msgData)
}

// ========= handler ============
type HandleResponse struct {
	Messages []string `json:"messages"`
	Logs []string `json:"log"`
	Data []byte `json:"data"`
}

func DoHandler(handlerFn func(deps *Extern, _env Env, msg []byte) HandleResponse, envPtr, msgPtr uint32) unsafe.Pointer{
	envData := Translate_range_custom(uintptr(envPtr))
	msgData := Translate_range_custom(uintptr(msgPtr))

	result := _do_handler(handlerFn, envData, msgData)

	data, err := ezjson.Marshal(result)
	if err != nil {
		return nil
	}

	return Package_message(data)
}

func _do_handler(handlerFn func(deps *Extern, _env Env, msg []byte) HandleResponse, envData, msgData []byte) HandleResponse{
	var env Env
	ezjson.Unmarshal(envData, &env)

	deps := make_dependencies()
	return handlerFn(&deps, env, msgData)
}

// =========== query ===================
type Binary []byte

func DoQuery(queryFn func(deps *Extern, msg []byte) Binary, msgPtr uint32) unsafe.Pointer{
	msgData := Translate_range_custom(uintptr(msgPtr))

	result := _do_query(queryFn, msgData)

	return Package_message(result)
}

func _do_query(handlerFn func(deps *Extern, msg []byte) Binary, msgData []byte) Binary {
	deps := make_dependencies()
	return handlerFn(&deps, msgData)
}
