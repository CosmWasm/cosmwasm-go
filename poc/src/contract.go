package src

import (
	"github.com/cosmwasm/cosmwasm-go/poc/std"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
)

func Init(deps *std.Extern, _env std.Env, msg []byte) std.CosmosResponse {
	return std.CosmosResponseDefault()
}

func Invoke(deps *std.Extern, _env std.Env, msg []byte) std.CosmosResponse {
	var handlerMsg HandleMsg
	err := ezjson.Unmarshal(msg, &handlerMsg)
	if err != nil {
		return std.CosmosResponse{
			Ok:  nil,
			Err: std.ToStdError(std.GenericErr{Msg: err.Error()}),
		}
	}

	if handlerMsg.Register != nil {
		return tryRegister(deps, _env, handlerMsg.Register)
	} else if handlerMsg.Sell != nil {
		return trySell(deps, _env, handlerMsg.Sell)
	}

	return std.CosmosResponse{
		Ok:  nil,
		Err: std.ToStdError(std.GenericErr{Msg: "unknowns function called!"}),
	}
}

func Query(deps *std.Extern, msg []byte) std.CosmosResponse {
	var queryMsg QueryMsg
	err := ezjson.Unmarshal(msg, &queryMsg)
	if err != nil {
		return std.CosmosResponse{
			Ok:  nil,
			Err: std.ToStdError(std.GenericErr{Msg: err.Error()}),
		}
	}

	if queryMsg.Get != nil {
		return tryGetDomain(deps, queryMsg.Get)
	}

	return std.CosmosResponse{
		Ok:  nil,
		Err: std.ToStdError(std.GenericErr{Msg: "unknowns function called!"}),
	}
}

func tryRegister(deps *std.Extern, _env std.Env, registerInfo *RegisterDomain) std.CosmosResponse {
	if err := deps.EStorage.Set([]byte(registerInfo.Domain), _env.Message.Sender); err != nil {
		return std.CosmosResponse{
			Ok:  nil,
			Err: std.ToStdError(std.GenericErr{Msg: err.Error()}),
		}
	}

	return std.CosmosResponseDefault()
}

func trySell(deps *std.Extern, _env std.Env, sellInfo *SellDomain) std.CosmosResponse {
	buyerCanonAddress, err := deps.EApi.CanonicalAddress(sellInfo.Buyer)
	if err != nil {
		return std.CosmosResponse{
			Ok:  nil,
			Err: std.ToStdError(std.GenericErr{Msg: err.Error()}),
		}
	}

	if err := deps.EStorage.Set([]byte(sellInfo.Domain), buyerCanonAddress); err != nil {
		return std.CosmosResponse{
			Ok:  nil,
			Err: std.ToStdError(std.GenericErr{Msg: err.Error()}),
		}
	}

	return std.CosmosResponseDefault()
}

func tryGetDomain(deps *std.Extern, queryInfo *GetOwner) std.CosmosResponse {
	owner, err := deps.EStorage.Get([]byte(queryInfo.Domain))
	if err != nil {
		return std.CosmosResponse{
			Ok:  nil,
			Err: std.ToStdError(std.GenericErr{Msg: err.Error()}),
		}
	}

	return std.CosmosResponse{
		Ok: &std.Result{
			Messages: []std.CosmosMsg{},
			Data:     string(owner),
			Log:      []std.LogAttribute{},
		},
		Err: nil,
	}
}
