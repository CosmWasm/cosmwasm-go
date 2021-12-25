package contractv2

import (
	"github.com/cosmwasm/cosmwasm-go/example/contractv2/src"
	"github.com/cosmwasm/cosmwasm-go/std"
	"unsafe"
)

func main() {}

//export execute
func execute(envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoExecute(src.Execute, envPtr, infoPtr, msgPtr)
}

//export query
func query(envPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoQuery(src.Query, envPtr, msgPtr)
}

//export instantiate
func instantiate(envPtr, infoPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoInstantiate(src.Instantiate, envPtr, infoPtr, msgPtr)
}

// export migrate
func migrate(envPtr, msgPtr uint32) unsafe.Pointer {
	return std.DoMigrate(src.Migrate, envPtr, msgPtr)
}
