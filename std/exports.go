//go:build cosmwasm
// +build cosmwasm

package std

import (
	"unsafe"

	"github.com/CosmWasm/cosmwasm-go/std/types"
)

type (
	// InstantiateFunc defines the function ran by contracts in instantiation.
	InstantiateFunc func(deps *Deps, env types.Env, messageInfo types.MessageInfo, messageBytes []byte) (*types.Response, error)
	// ExecuteFunc defines the function ran by contracts in message execution.
	ExecuteFunc func(deps *Deps, env types.Env, messageInfo types.MessageInfo, messageBytes []byte) (*types.Response, error)
	// MigrateFunc defines the function ran by contracts in migration.
	MigrateFunc func(deps *Deps, env types.Env, messageBytes []byte) (*types.Response, error)
	// SudoFunc defines the function ran by contracts in sudo message execution.
	SudoFunc func(deps *Deps, env types.Env, messageBytes []byte) (*types.Response, error)
	// QueryFunc defines the function ran by the contracts in query execution.
	QueryFunc func(deps *Deps, env types.Env, messageBytes []byte) ([]byte, error)
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
// and executes the contract's message execution logic.
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

// DoMigrate converts the environment and message pointers to concrete golang objects
// and execute the contract migration logic.
func DoMigrate(migrateFunc MigrateFunc, envPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	resp, err := migrateFunc(&deps, env, msgData)
	if resp == nil {
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

// DoSudo converts the environment and message pointers to concrete golang objects
// and executes the contract's sudo message execution logic.
func DoSudo(sudoFunc SudoFunc, envPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	resp, err := sudoFunc(&deps, env, msgData)
	if resp == nil {
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

// DoQuery converts the environment and info pointers to concrete golang objects
// and executes the contract's query logic.
func DoQuery(queryFunc QueryFunc, envPtr, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err)
	}

	deps := make_dependencies()
	respBytes, err := queryFunc(&deps, env, msgData)
	if err != nil {
		return StdErrResult(err)
	}

	result := &types.QueryResponse{
		Ok: respBytes,
	}
	data, err := result.MarshalJSON()
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

//export interface_version_8
func interface_version_8() {}
