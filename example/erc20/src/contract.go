package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func Init(deps *std.Deps, env std.Env, info std.MessageInfo, msg []byte) (*std.InitResultOk, error) {
	initMsg := InitMsg{}

	ownerShip := NewOwnership(deps)
	err := ezjson.Unmarshal(msg, &initMsg)
	if err != nil {
		return nil, err
	}
	if err = initMsg.Validate(); err != nil {
		return nil, err
	}
	deps.Api.Debug("*** Init Called ***")

	erc20Protocol := NewErc20Protocol(State{
		NameOfToken:   initMsg.Name,
		SymbolOfToken: initMsg.Symbol,
		DecOfTokens:   initMsg.Decimal,
		TotalSupplyOf: initMsg.TotalSupply,
	}, deps, &info)

	owner, err := deps.Api.CanonicalAddress(info.Sender)
	if err != nil {
		return nil, err
	}
	ownerShip.Owned(owner)

	erc20Protocol.Assign(owner, 10000)
	//saving state and owner info
	erc20Protocol.SaveState()
	ownerShip.SaveOwner()
	return &std.InitResultOk{
		Ok: std.InitResponse{
			Attributes: []std.Attribute{{"hello", "world"}},
		},
	}, nil
}

func Invoke(deps *std.Deps, env std.Env, info std.MessageInfo, msg []byte) (*std.HandleResultOk, error) {
	return handleInvokeMessage(deps, env, info, msg)
}

func Query(deps *std.Deps, env std.Env, msg []byte) (*std.QueryResponse, error) {
	return handleQuery(deps, env, msg)
}
