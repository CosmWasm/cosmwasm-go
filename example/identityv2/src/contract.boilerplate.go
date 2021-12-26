package src

import (
	std "github.com/cosmwasm/cosmwasm-go/std"
	types "github.com/cosmwasm/cosmwasm-go/std/types"
)

// ExecuteMsg is the union type used to process execution messages towards the contract.
type ExecuteMsg struct {
	CreateIdentity *MsgCreateIdentity `json:"create_identity"`
	DeleteIdentity *MsgDelete         `json:"delete_identity"`
	UpdateCity     *MsgUpdateCity     `json:"update_city"`
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, messageBytes []byte) (*types.Response, error) {
	msg := new(ExecuteMsg)
	err := msg.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}
	switch {
	case msg.CreateIdentity != nil:
		resp, err := Contract{}.ExecCreateIdentity(deps, &env, &info, msg.CreateIdentity)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case msg.DeleteIdentity != nil:
		resp, err := Contract{}.ExecDeleteIdentity(deps, &env, &info, msg.DeleteIdentity)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case msg.UpdateCity != nil:
		resp, err := Contract{}.ExecUpdateCity(deps, &env, &info, msg.UpdateCity)
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		panic(1)
	}
}

func (x *MsgCreateIdentity) AsExecuteMsg() ExecuteMsg {
	return ExecuteMsg{CreateIdentity: x}
}

func (x *MsgDelete) AsExecuteMsg() ExecuteMsg {
	return ExecuteMsg{DeleteIdentity: x}
}

func (x *MsgUpdateCity) AsExecuteMsg() ExecuteMsg {
	return ExecuteMsg{UpdateCity: x}
}

// QueryMsg is the union type used to process queries towards the contract.
type QueryMsg struct {
	Identity *QueryIdentity `json:"identity"`
}

func Query(deps *std.Deps, env types.Env, queryBytes []byte) ([]byte, error) {
	query := new(QueryMsg)
	err := query.UnmarshalJSON(queryBytes)
	if err != nil {
		return nil, err
	}
	switch {
	case query.Identity != nil:
		resp, err := Contract{}.QueryIdentity(deps, &env, query.Identity)
		if err != nil {
			return nil, err
		}
		return resp.MarshalJSON()
	default:
		panic(1)
	}
}

func (x *QueryIdentity) AsQueryMsg() QueryMsg {
	return QueryMsg{Identity: x}
}

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, instantiateBytes []byte) (*types.Response, error) {
	return &types.Response{}, nil
}
func Migrate(deps *std.Deps, env types.Env, migrateBytes []byte) (*types.Response, error) {
	return &types.Response{}, nil
}
