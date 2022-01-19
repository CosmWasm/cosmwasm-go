package src

import (
	xd "github.com/cosmwasm/cosmwasm-go/example/identityv2/src/imp"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

type Contract struct{}

// CreateIdentity creates a new identity.
// +cw:exec
func (c Contract) CreateIdentity(deps *std.Deps, _ *types.Env, info *types.MessageInfo, msg *MsgCreateIdentity) (*types.Response, error) {
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

// UpdateCity updates a Person's city.
// +cw:exec
func (c Contract) UpdateCity(deps *std.Deps, _ *types.Env, info *types.MessageInfo, msg *MsgUpdateCity) (*types.Response, error) {
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

// DeleteIdentity deletes the Person object associated with the sender.
// +cw:exec
func (c Contract) DeleteIdentity(deps *std.Deps, _ *types.Env, info *types.MessageInfo, _ *MsgDelete) (*types.Response, error) {
	err := PersonState.Delete(deps.Storage, Person{Address: info.Sender})
	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}

// ImportedMessage does nothing. It tests codegen on imported messages.
// +cw:exec
func (c Contract) ImportedMessage(_ *std.Deps, _ *types.Env, _ *types.MessageInfo, _ *xd.ImportedMessage) (*types.Response, error) {
	return &types.Response{}, nil
}

// Migrate does nothing. It tests codegen on migration methods.
// +cw:migrate
func (c Contract) Migrate(_ *std.Deps, _ *types.Env, _ *MsgMigrate) (*types.Response, error) {
	return &types.Response{}, nil
}

// Instantiate does nothing. It tests codegen on instantiate methods.
// +cw:instantiate
func (c Contract) Instantiate(_ *std.Deps, _ *types.Env, _ *types.MessageInfo, _ *MsgInstantiate) (*types.Response, error) {
	return &types.Response{}, nil
}

// QueryIdentity queries a Person.
// +cw:query
func (c Contract) QueryIdentity(deps *std.Deps, _ *types.Env, msg *QueryIdentity) (*QueryIdentityResponse, error) {
	res, err := PersonState.Read(deps.Storage, msg.ID)
	if err != nil {
		return nil, err
	}

	return &QueryIdentityResponse{Person: &res}, nil
}

// QueryImported does nothing. It tests codegen on imported query messages.
// +cw:query
func (c Contract) QueryImported(_ *std.Deps, _ *types.Env, _ *xd.ImportedQuery) (*xd.ImportedQueryResponse, error) {
	return &xd.ImportedQueryResponse{}, nil
}
