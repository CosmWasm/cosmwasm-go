package src

import (
	xd "github.com/cosmwasm/cosmwasm-go/example/identityv2/src/imp"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

type Contract struct{}

// +cw:exec
func (c Contract) CreateIdentity(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *MsgCreateIdentity) (*types.Response, error) {
	panic("impl")
}

// +cw:exec
func (c Contract) ImportedMessage(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *xd.ImportedMessage) (*types.Response, error) {
	panic("impl")
}

// +cw:migrate
func (c Contract) Migrate(deps *std.Deps, env *types.Env, msg *MsgMigrate) (*types.Response, error) {
	panic("impl")
}

// +cw:instantiate
func (c Contract) Instantiate(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *MsgInstantiate) (*types.Response, error) {
	panic("impl")
}

// +cw:query
func (c Contract) QueryIdentity(deps *std.Deps, env *types.Env, msg *QueryIdentity) (*QueryIdentityResponse, error) {
	panic("impl")
}

// +cw:query
func (c Contract) QueryImported(deps *std.Deps, env *types.Env, msg *xd.ImportedQuery) (*xd.ImportedQueryResponse, error) {
	panic("impl")
}
