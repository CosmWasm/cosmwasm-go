package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func Init(deps *std.Extern, _env std.Env, msg []byte) (*std.CosmosResponseOk, std.CosmosResponseError) {
	return &std.CosmosResponseOk{
		Ok: std.Result{
			Messages: nil,
			Log: []std.LogAttribute{
				{Key: "Key1", Value: "Value1"},
				{Key: "Key2", Value: "Value2"},
			},
		},
	}, ""
}

func Invoke(deps *std.Extern, _env std.Env, msg []byte) (*std.CosmosResponseOk, std.CosmosResponseError) {
	var handlerMsg HandleMsg
	err := ezjson.Unmarshal(msg, &handlerMsg)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Testing generic error result", "")
	}

	return nil, std.GenerateError(std.GenericError, "Testing generic error result", "")
}

func Query(deps *std.Extern, msg []byte) (*std.CosmosResponseOk, std.CosmosResponseError) {
	var queryMsg QueryMsg
	err := ezjson.Unmarshal(msg, &queryMsg)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Testing generic error result", "")
	}

	return nil, std.GenerateError(std.GenericError, "Testing generic error result", "")
}
