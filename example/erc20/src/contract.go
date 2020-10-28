package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func Init(deps *std.Extern, env std.Env, info std.MessageInfo, msg []byte) (*std.InitResultOk, *std.CosmosResponseError) {
	initMsg := InitMsg{}
	ownerShip := NewOwnership(deps)
	e := ezjson.Unmarshal(msg, &initMsg)
	if e != nil {
		return nil, std.GenerateError(std.GenericError, "Unmarshal initMsg failed "+e.Error(), "")
	}
	if err := initMsg.Validate(); err != nil {
		return nil, err
	}

	erc20Protocol := NewErc20Protocol(State{
		NameOfToken:   initMsg.Name,
		SymbolOfToken: initMsg.Symbol,
		DecOfTokens:   initMsg.Decimal,
		TotalSupplyOf: initMsg.TotalSupply,
	}, deps, &info)

	owner, err := deps.EApi.CanonicalAddress(info.Sender)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Invalid Sender: "+err.Error(), "")
	}
	ownerShip.Owned(owner)

	erc20Protocol.Assign(owner, 10000)
	//saving state and owner info
	erc20Protocol.SaveState()
	ownerShip.SaveOwner()
	return std.InitResultOkOkDefault(), nil
}

func Invoke(deps *std.Extern, env std.Env, info std.MessageInfo, msg []byte) (*std.HandleResultOk, *std.CosmosResponseError) {
	return handleInvokeMessage(deps, env, info, msg)
}

func Query(deps *std.Extern, env std.Env, msg []byte) (*std.QueryResponseOk, *std.CosmosResponseError) {
	return handleQuery(deps, env, msg)
}
