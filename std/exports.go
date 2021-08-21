// +build cosmwasm

package std

import (
	"strings"
	"unsafe"
)

type ContractError struct {
	Err string `json:"error"`
}

func StdErrResult(err error, prefix string) unsafe.Pointer {
	raw := err.Error()
	if prefix != "" {
		raw = prefix + ": " + raw
	}
	clean := strings.Replace(raw, `"`, `\"`, -1)
	msg := `{"error":"` + clean + `"}`
	return Package_message([]byte(msg))
}

func make_dependencies() Deps {
	return Deps{
		Storage: ExternalStorage{},
		Api:     ExternalApi{},
		Querier: ExternalQuerier{},
	}
}

func parseInfo(infoPtr uint32) (MessageInfo, error) {
	infoData := TranslateToSlice(uintptr(infoPtr))
	var info MessageInfo
	err := info.UnmarshalJSON(infoData)
	return info, err
}

// ========== init ==============
func DoInit(initFn func(*Deps, Env, MessageInfo, []byte) (*InitResultOk, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse Env")
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, err := initFn(&deps, env, info, msgData)
	if ok == nil || err != nil {
		return StdErrResult(err, "Init")
	}

	data, err := ok.MarshalJSON()
	if err != nil {
		return StdErrResult(err, "Marshal Response")
	}
	return Package_message(data)
}

// ========= handler ============
func DoHandler(handlerFn func(*Deps, Env, MessageInfo, []byte) (*HandleResultOk, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse Env")
	}

	info, err := parseInfo(infoPtr)
	if err != nil {
		return StdErrResult(err, "Parse Info")
	}

	deps := make_dependencies()
	msgData := Translate_range_custom(uintptr(msgPtr))
	ok, err := handlerFn(&deps, env, info, msgData)
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
func DoMigrate(migrateFn func(*Deps, Env, MessageInfo, []byte) (*MigrateResultOk, error), envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse Env")
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
func DoQuery(queryFn func(*Deps, Env, []byte) (*QueryResponse, error), envPtr, msgPtr uint32) unsafe.Pointer {
	msgData := Translate_range_custom(uintptr(msgPtr))
	env := Env{}
	envData := TranslateToSlice(uintptr(envPtr))
	err := env.UnmarshalJSON(envData)
	if err != nil {
		return StdErrResult(err, "Parse Env")
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
