package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

//go:generate ../../../bin/generator contract contract.go Contract
type Contract struct{}

func (c Contract) ExecCreateIdentity(deps *std.Deps, _ *types.Env, info *types.MessageInfo, msg *MsgCreateIdentity) (*types.Response, error) {
	err := PersonState.Create(deps.Storage, Person{
		Address:    info.Sender,
		Name:       msg.Name,
		Surname:    msg.Surname,
		City:       msg.City,
		PostalCode: msg.PostalCode,
	})

	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}

func (c Contract) ExecUpdateCity(deps *std.Deps, _ *types.Env, info *types.MessageInfo, msg *MsgUpdateCity) (*types.Response, error) {
	err := PersonState.Update(deps.Storage, info.Sender, func(old *Person) (*Person, error) {
		if old == nil {
			return nil, ErrPersonNotFound
		}

		old.City = msg.City
		old.PostalCode = msg.PostalCode

		return old, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}

func (c Contract) ExecDeleteIdentity(deps *std.Deps, _ *types.Env, info *types.MessageInfo, _ *MsgDelete) (*types.Response, error) {
	err := PersonState.Delete(deps.Storage, Person{Address: info.Sender})
	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}

func (c Contract) QueryIdentity(deps *std.Deps, _ *types.Env, query *QueryIdentity) (*Person, error) {
	person, err := PersonState.Read(deps.Storage, query.ID)
	if err != nil {
		return nil, err
	}
	return &person, nil
}
