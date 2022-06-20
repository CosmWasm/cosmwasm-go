package main

import (
	"unsafe"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src"
	"github.com/CosmWasm/cosmwasm-go/std"
)

// main is not used as a WASM contract works with callbacks only.
func main() {}

// instantiate is a WASM contract entrypoint.
//export instantiate
func instantiate(envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoInstantiate(src.Instantiate, envPtr, infoPtr, msgPtr)
}

// execute is a WASM contract entrypoint.
//export execute
func execute(envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoExecute(src.Execute, envPtr, infoPtr, msgPtr)
}

// migrate is a WASM contract entrypoint.
//export migrate
func migrate(envPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoMigrate(src.Migrate, envPtr, msgPtr)
}

// sudo is a WASM contract entrypoint.
//export sudo
func sudo(envPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoSudo(src.Sudo, envPtr, msgPtr)
}

// reply is a WASM contract entrypoint.
//export reply
func reply(envPtr, replyPtr uint32) unsafe.Pointer {
	return std.DoReply(src.Reply, envPtr, replyPtr)
}

// query is a WASM contract entrypoint.
//export query
func query(envPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoQuery(src.Query, envPtr, msgPtr)
}
