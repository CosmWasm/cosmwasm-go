// +build cosmwasm

package std

import (
	"unsafe"

	"github.com/cosmwasm/cosmwasm-go/std/types"
)

func StdErrResult(err error, _prefix string) unsafe.Pointer {
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

// ========== instantiate ==============
func DoInstantiate(instantiateFn func(*Deps, types.Env, types.MessageInfo, []byte) (*types.ContractResult, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse types.Env")
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, err := instantiateFn(&deps, env, info, msgData)
	if ok == nil || err != nil {
		return StdErrResult(err, "Instantiate")
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err, "Marshal Response")
	}
	return Package_message(data)
}

// ========= execute ============
func DoExecute(executeFn func(*Deps, types.Env, types.MessageInfo, []byte) (*types.ContractResult, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse types.Env")
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, err := executeFn(&deps, env, info, msgData)
	if ok == nil || err != nil {
		return StdErrResult(err, "Handle")
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err, "Marshal Response")
	}
	return Package_message(data)
}

// ========= migrate ============
func DoMigrate(migrateFn func(*Deps, types.Env, types.MessageInfo, []byte) (*types.ContractResult, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := types.Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse types.Env")
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, err := migrateFn(&deps, env, info, msgData)
	if ok == nil {
		return StdErrResult(err, "Migrate")
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err, "Marshal Response")
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
		return StdErrResult(err, "Parse types.Env")
	}

	deps := make_dependencies()
	ok, err := queryFn(&deps, env, msgData)
	if ok == nil {
		return StdErrResult(err, "Query")
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err, "Marshal Response")
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
