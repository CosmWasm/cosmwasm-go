package src

import (
	std "github.com/cosmwasm/cosmwasm-go/std"
	types "github.com/cosmwasm/cosmwasm-go/std/types"
)

// ExecuteMsg is the union type used to process execution messages towards the contract.
type ExecuteMsg struct {
	Echo *MsgEcho `json:"echo"`
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, messageBytes []byte) (*types.Response, error) {
	msg := new(ExecuteMsg)
	err := msg.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}
	switch {
	case msg.Echo != nil:
		resp, err := Contract{}.ExecEcho(deps, &env, &info, msg.Echo)
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		panic(1)
	}
}

// QueryMsg is the union type used to process queries towards the contract.
type QueryMsg struct {
	Key *QueryKey `json:"key"`
}

func Query(deps *std.Deps, env types.Env, queryBytes []byte) ([]byte, error) {
	query := new(QueryMsg)
	err := query.UnmarshalJSON(queryBytes)
	if err != nil {
		return nil, err
	}
	switch {
	case query.Key != nil:
		resp, err := Contract{}.QueryKey(deps, &env, query.Key)
		if err != nil {
			return nil, err
		}
		return resp.MarshalJSON()
	default:
		panic(1)
	}
}

func (x *QueryKey) AsQueryMsg() *QueryMsg {
	return &QueryMsg{Key: x}
}

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, instantiateBytes []byte) (*types.Response, error) {
	initMsg := new(MsgInit)
	err := initMsg.UnmarshalJSON(instantiateBytes)
	if err != nil {
		return nil, err
	}
	return Contract{}.Instantiate(deps, &env, &info, initMsg)
}

func Migrate(deps *std.Deps, env types.Env, migrateBytes []byte) (*types.Response, error) {
	msg := new(MsgMigrate)
	err := msg.UnmarshalJSON(migrateBytes)
	if err != nil {
		return nil, err
	}
	return Contract{}.Migrate(deps, &env, msg)
}
