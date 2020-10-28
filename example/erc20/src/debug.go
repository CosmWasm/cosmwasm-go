// +build cosmwasm

package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func init() {
	ezjson.SetDisplay(std.Wasmlog)
}
