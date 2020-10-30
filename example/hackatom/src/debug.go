// +build cosmwasm

package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// set this to true if you want to do deep debuging of ezjson parsing inside the wasm blog
// (using debug import to print)
const DEBUG_EZJSON = true

func init() {
	if DEBUG_EZJSON {
		ezjson.SetDisplay(std.Wasmlog)
	}
}
