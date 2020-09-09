package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
)

func Init(deps *std.Extern, _env std.Env, msg []byte) (*std.CosmosResponseOk, *std.CosmosResponseError) {
	return nil, std.GenerateError(std.GenericError, "unimplemented", "")
}

func Invoke(deps *std.Extern, _env std.Env, msg []byte) (*std.CosmosResponseOk, *std.CosmosResponseError) {
	return nil, std.GenerateError(std.GenericError, "unimplemented", "")
}

func Query(deps *std.Extern, msg []byte) (*std.CosmosResponseOk, *std.CosmosResponseError) {
	return nil, std.GenerateError(std.GenericError, "unimplemented", "")
}
