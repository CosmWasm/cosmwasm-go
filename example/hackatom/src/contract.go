package src

import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func Init(deps *std.Extern, env std.Env, info std.MessageInfo, msg []byte) (*std.InitResultOk, error) {
	initMsg := InitMsg{}
	err := ezjson.Unmarshal(msg, &initMsg)
	if err != nil {
		return nil, err
	}

	owner, err := deps.EApi.CanonicalAddress(info.Sender)
	if err != nil {
		return nil, err
	}
	state := State{
		Count: initMsg.Count,
		Owner: owner,
	}

	bz, err := ezjson.Marshal(state)
	if err != nil {
		return nil, err
	}
	// TODO: this should not return error
	// TODO: change names EApi -> Api
	err = deps.EStorage.Set(StateKey, bz)
	if err != nil {
		return nil, err
	}

	return &std.InitResultOk{}, nil
}

func Handle(deps *std.Extern, env std.Env, info std.MessageInfo, msg []byte) (*std.HandleResultOk, error) {
	return nil, errors.New("Not implemented")
}

func Query(deps *std.Extern, env std.Env, msg []byte) (*std.QueryResponseOk, error) {
	return nil, errors.New("Not implemented")
}
