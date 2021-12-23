package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

type MsgEcho struct {
	Msg string
}

func (m *MsgEcho) UnmarshalJSON(b []byte) error {
	panic("")
}

type QueryKey struct {
	Key []byte
}

func (x *QueryKey) UnmarshalJSON(b []byte) error {
	panic("impl")
}

type QueryKeyResponse struct {
	Value []byte
}

func (x QueryKeyResponse) MarshalJSON() ([]byte, error) {
	panic("impl")
}

type Contract struct{}

func (c Contract) ExecEcho(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *MsgEcho) (*types.Response, error) {
	// do stuff
	return &types.Response{}, nil
}

func (c Contract) QueryKey(deps *std.Deps, env *types.Env, query *QueryKey) (*QueryKeyResponse, error) {
	return &QueryKeyResponse{Value: deps.Storage.Get(query.Key)}, nil
}
