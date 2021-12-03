package src

import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

var (
	ErrUnknownRequest = errors.New("unknown request")
)

func Instantiate(_ *std.Deps, _ types.Env, _ types.MessageInfo, _ []byte) (*types.Response, error) {
	return &types.Response{}, nil
}

func Execute(deps *std.Deps, _ types.Env, messageInfo types.MessageInfo, messageBytes []byte) (*types.Response, error) {
	msg := new(ExecuteMsg)
	err := msg.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}

	switch {
	case msg.CreateIdentity != nil:
		return executeCreateIdentity(deps, messageInfo, msg.CreateIdentity)
	case msg.UpdateCity != nil:
		return executeUpdateCity(deps, messageInfo, msg.UpdateCity)
	case msg.DeleteIdentity != nil:
		return executeDeleteIdentity(deps, messageInfo, msg.DeleteIdentity)
	default:
		return nil, ErrUnknownRequest
	}
}

func executeDeleteIdentity(deps *std.Deps, info types.MessageInfo, _ *MsgDelete) (*types.Response, error) {
	err := PersonState.Delete(deps.Storage, Person{Address: info.Sender})
	if err != nil {
		return nil, err
	}
	return &types.Response{}, nil
}

func executeUpdateCity(deps *std.Deps, info types.MessageInfo, update *MsgUpdateCity) (*types.Response, error) {
	err := PersonState.Update(deps.Storage, info.Sender, func(old *Person) (*Person, error) {
		if old == nil {
			return nil, ErrPersonNotFound
		}
		old.City = update.City
		old.PostalCode = update.PostalCode
		return old, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.Response{}, err
}

func executeCreateIdentity(deps *std.Deps, info types.MessageInfo, identity *MsgCreateIdentity) (*types.Response, error) {
	err := PersonState.Create(deps.Storage, Person{
		Address:    info.Sender,
		Name:       identity.Name,
		Surname:    identity.Surname,
		City:       identity.City,
		PostalCode: identity.PostalCode,
	})
	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}

func Query(deps *std.Deps, env types.Env, messageBytes []byte) ([]byte, error) {
	q := new(QueryMsg)
	err := q.UnmarshalJSON(messageBytes)
	if err != nil {
		return nil, err
	}

	switch {
	case q.Identity != nil:
		return queryIdentity(deps, q.Identity)
	default:
		return nil, ErrUnknownRequest
	}
}

func queryIdentity(deps *std.Deps, identity *QueryIdentity) ([]byte, error) {
	person, err := PersonState.Read(deps.Storage, identity.ID)
	if err != nil {
		return nil, err
	}

	return person.MarshalJSON()
}

func Migrate(_ *std.Deps, _ types.Env, _ []byte) (*types.Response, error) {
	return &types.Response{}, nil
}
