package src

import (
	errors "errors"
	imp "github.com/cosmwasm/cosmwasm-go/example/identityv2/src/imp"
	std "github.com/cosmwasm/cosmwasm-go/std"
	types "github.com/cosmwasm/cosmwasm-go/std/types"
)

type QueryMsg struct {
	QueryImported *imp.ImportedQuery `json:"query_imported"`
	QueryIdentity *QueryIdentity     `json:"query_identity"`
}

type ExecuteMsg struct {
	DeleteIdentity  *MsgDelete           `json:"delete_identity"`
	ImportedMessage *imp.ImportedMessage `json:"imported_message"`
	UpdateCity      *MsgUpdateCity       `json:"update_city"`
	CreateIdentity  *MsgCreateIdentity   `json:"create_identity"`
}

func (x *QueryIdentity) AsQueryMsg() QueryMsg {
	return QueryMsg{QueryIdentity: x}
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

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, messageBytes []byte) (*types.Response, error) {
	msg := new(ExecuteMsg)
	err := msg.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}
	switch {
	case msg.CreateIdentity != nil:
		return Contract{}.CreateIdentity(deps, &env, &info, msg.CreateIdentity)
	case msg.DeleteIdentity != nil:
		return Contract{}.DeleteIdentity(deps, &env, &info, msg.DeleteIdentity)
	case msg.ImportedMessage != nil:
		return Contract{}.ImportedMessage(deps, &env, &info, msg.ImportedMessage)
	case msg.UpdateCity != nil:
		return Contract{}.UpdateCity(deps, &env, &info, msg.UpdateCity)
	default:
		return nil, errors.New("unknown request")
	}
}

func Query(deps *std.Deps, env types.Env, queryBytes []byte) ([]byte, error) {
	query := new(QueryMsg)
	err := query.UnmarshalJSON(queryBytes)
	if err != nil {
		return nil, err
	}
	switch {
	case query.QueryIdentity != nil:
		resp, err := Contract{}.QueryIdentity(deps, &env, query.QueryIdentity)
		if err != nil {
			return nil, err
		}
		return resp.MarshalJSON()
	case query.QueryImported != nil:
		resp, err := Contract{}.QueryImported(deps, &env, query.QueryImported)
		if err != nil {
			return nil, err
		}
		return resp.MarshalJSON()
	default:
		return nil, errors.New("unknown request")
	}
}

func Migrate(deps *std.Deps, env types.Env, messageBytes []byte) (*types.Response, error) {
	msg := new(MsgMigrate)
	err := msg.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}
	return Contract{}.Migrate(deps, &env, msg)
}

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, messageBytes []byte) (*types.Response, error) {
	msg := new(MsgInstantiate)
	err := msg.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}
	return Contract{}.Instantiate(deps, &env, &info, msg)
}
