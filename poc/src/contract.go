package src

import (
	"github.com/cosmwasm/cosmwasm-go/poc/std"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
)

func Init(deps *std.Extern, _env std.Env, msg []byte) std.InitResponse {
	return std.InitResponse{}
}

func Invoke(deps *std.Extern, _env std.Env, msg []byte) std.HandleResponse {
	var handlerMsg HandleMsg
	err := ezjson.Unmarshal(msg, &handlerMsg)
	if err != nil {
		return std.HandleResponse{
			Messages: nil,
			Logs:     []string{err.Error()},
			Data:     nil,
		}
	}

	if handlerMsg.Register != nil {
		return tryRegister(deps, _env, handlerMsg.Register)
	} else if handlerMsg.Sell != nil {
		return trySell(deps, _env, handlerMsg.Sell)
	}

	return std.HandleResponse{
		Messages: nil,
		Logs:     []string{"unknowns function called!"},
		Data:     nil,
	}
}

func Query(deps *std.Extern, msg []byte) []byte {
	var queryMsg QueryMsg
	err := ezjson.Unmarshal(msg, &queryMsg)
	if err != nil {
		return []byte(err.Error())
	}

	if queryMsg.Get != nil {
		return tryGetDomain(deps, queryMsg.Get)
	}

	return []byte("unknowns function called!")
}

func tryRegister(deps *std.Extern, _env std.Env, registerInfo *RegisterDomain) std.HandleResponse {
	if err := deps.EStorage.Set([]byte(registerInfo.Domain), _env.Message.Sender); err != nil {
		return std.HandleResponse{
			Messages: nil,
			Logs:     []string{err.Error()},
			Data:     nil,
		}
	}

	return std.HandleResponse{}
}

func trySell(deps *std.Extern, _env std.Env, sellInfo *SellDomain) std.HandleResponse {
	buyerCanonAddress, err := deps.EApi.CanonicalAddress(sellInfo.Buyer)
	if err != nil {
		return std.HandleResponse{
			Messages: nil,
			Logs:     []string{err.Error()},
			Data:     nil,
		}
	}

	if err := deps.EStorage.Set([]byte(sellInfo.Domain), buyerCanonAddress); err != nil {
		return std.HandleResponse{
			Messages: nil,
			Logs:     []string{err.Error()},
			Data:     nil,
		}
	}

	return std.HandleResponse{}
}

func tryGetDomain(deps *std.Extern, queryInfo *GetOwner) []byte {
	owner, err := deps.EStorage.Get([]byte(queryInfo.Domain))
	if err != nil {
		return []byte(err.Error())
	}

	return owner
}
