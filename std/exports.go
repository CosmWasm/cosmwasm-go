//go:build cosmwasm
// +build cosmwasm

package std

import (
	"unsafe"

	"github.com/cosmwasm/cosmwasm-go/std/types"
)

type (
	// InstantiateFunc defines the function ran by contracts in instantiation.
	InstantiateFunc func(deps *Deps, env types.Env, messageInfo types.MessageInfo, messageBytes []byte) (*types.Response, error)
	// ExecuteFunc defines the function ran by contracts in message execution.
	ExecuteFunc func(deps *Deps, env types.Env, messageInfo types.MessageInfo, messageBytes []byte) (*types.Response, error)
)

func StdErrResult(err error) unsafe.Pointer {
	wrapped := types.ContractResult{Err: err.Error()}
	bz, _ := wrapped.MarshalJSON()
	return Package_message(bz)
}

func make_dependencies() Deps {
	return Deps{
		Storage: ExternalStorage{},
		Api:     ExternalApi{},
		Querier: ExternalQuerier{},
	}
}

func parseInfo(infoPtr uint32) (types.MessageInfo, error) {
	infoData := TranslateToSlice(uintptr(infoPtr))
	var info types.MessageInfo
	err := info.UnmarshalJSON(infoData)
	return info, err
}

// DoInstantiate converts the environment, info and message pointers to concrete golang objects
// and executes the contract's instantiation function, returning a reference of the result.
func DoInstantiate(instantiateFunc InstantiateFunc, envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	resp, err := instantiateFunc(&deps, env, info, msgData)
	if err != nil {
		return StdErrResult(err)
	}

	result := &types.ContractResult{
		Ok: resp,
	}
	data, err := result.MarshalJSON()
	if err != nil {
		return StdErrResult(err)
	}
	return Package_message(data)
}

// DoExecute converts the environment, info and message pointers to concrete golang objects
// and executes contract's message execution logic.
func DoExecute(executeFunc ExecuteFunc, envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	resp, err := executeFunc(&deps, env, info, msgData)
	if err != nil {
		return StdErrResult(err)
	}

	result := &types.ContractResult{Ok: resp}

	data, err := result.MarshalJSON()
	if err != nil {
		return StdErrResult(err)
	}
	return Package_message(data)
}

// ========= migrate ============
func DoMigrate(migrateFn func(*Deps, types.Env, types.MessageInfo, []byte) (*types.ContractResult, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, err := migrateFn(&deps, env, info, msgData)
	if ok == nil {
		return StdErrResult(err)
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err)
	}
	return Package_message(data)
}

// =========== query ===================
func DoQuery(queryFn func(*Deps, types.Env, []byte) (*types.QueryResponse, error), envPtr, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	ok, err := queryFn(&deps, env, msgData)
	if ok == nil {
		return StdErrResult(err)
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err)
	}
	return Package_message(data)
}

//export allocate
func allocate(size uint32) unsafe.Pointer {
	ptr, _ := Build_region(size, 0)
	return ptr
}

//export deallocate
func deallocate(pointer unsafe.Pointer) {
	Deallocate(pointer)
}

//export interface_version_7
func interface_version_7() {}
