package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
)

func Init(deps *std.Extern, _env std.Env, msg []byte) (*std.CosmosResponseOk, *std.CosmosResponseError) {
	tester := Tester{deps: deps}
	e := tester.DoTest()
	if e == nil {
		std.DisplayMessage([]byte("test success"))
	} else {
		std.DisplayMessage([]byte("Test failed + " + e.Error()))
	}
	return &std.CosmosResponseOk{
		Ok: std.Result{
			Messages: nil,
			Log: []std.LogAttribute{
				{Key: "Key1", Value: "Value1"},
				{Key: "Key2", Value: "Value2"},
			},
		},
	}, nil
}

func Invoke(deps *std.Extern, _env std.Env, msg []byte) (*std.CosmosResponseOk, *std.CosmosResponseError) {
	return nil, std.GenerateError(std.GenericError, "unimplemented", "")
}

func Query(deps *std.Extern, msg []byte) (*std.CosmosResponseOk, *std.CosmosResponseError) {
	return nil, std.GenerateError(std.GenericError, "unimplemented", "")
}
