package src

import (
	"errors"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func Init(deps *std.Deps, env std.Env, info std.MessageInfo, msg []byte) (*std.InitResultOk, error) {
	initMsg := InitMsg{}
	err := ezjson.Unmarshal(msg, &initMsg)
	if err != nil {
		return nil, err
	}

	state := State{
		Count: initMsg.Count,
		Owner: info.Sender,
	}

	err = SaveState(deps.Storage, &state)
	if err != nil {
		return nil, err
	}
	return &std.InitResultOk{}, nil
}

func Handle(deps *std.Deps, env std.Env, info std.MessageInfo, data []byte) (*std.HandleResultOk, error) {
	msg := HandleMsg{}
	err := ezjson.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.Increment.Delta != 0:
		return handleIncrement(deps, &env, &info, msg.Increment)
	case msg.Reset.Value != 0:
		return handleReset(deps, &env, &info, msg.Reset)
	default:
		return nil, errors.New("Unknown HandleMsg")
	}
}

func handleIncrement(deps *std.Deps, env *std.Env, info *std.MessageInfo, msg Increment) (*std.HandleResultOk, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	state.Count += msg.Delta

	err = SaveState(deps.Storage, state)
	if err != nil {
		return nil, err
	}
	return &std.HandleResultOk{}, nil
}

func handleReset(deps *std.Deps, env *std.Env, info *std.MessageInfo, msg Reset) (*std.HandleResultOk, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if info.Sender != state.Owner {
		return nil, errors.New("Unauthorized")
	}
	state.Count = msg.Value

	err = SaveState(deps.Storage, state)
	if err != nil {
		return nil, err
	}
	return &std.HandleResultOk{}, nil
}

func Query(deps *std.Deps, env std.Env, data []byte) (*std.QueryResponseOk, error) {
	msg := QueryMsg{}
	err := ezjson.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.Count.A != "":
		// ignore this A field, it is just a placeholder for serialization
		return queryCount(deps, &env)
	default:
		return nil, errors.New("Unknown QueryMsg")
	}
}

func queryCount(deps *std.Deps, env *std.Env) (*std.QueryResponseOk, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	res := CountResponse{
		Count: state.Count,
	}
	bz, err := ezjson.Marshal(res)
	if err != nil {
		return nil, err
	}

	return &std.QueryResponseOk{
		Ok: bz,
	}, nil
}
