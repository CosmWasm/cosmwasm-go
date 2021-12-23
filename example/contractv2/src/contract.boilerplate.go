package src

import (
	std "github.com/cosmwasm/cosmwasm-go/std"
	types "github.com/cosmwasm/cosmwasm-go/std/types"
)

// QueryMsg is the union type used to process queries towards the contract.
type QueryMsg struct {
	Key *QueryKey `json:"key"`
}

func (x *QueryMsg) UnmarshalJSON(b []byte) error {
	panic(0)
}
func (x *QueryMsg) MarshalJSON() ([]byte, error) {
	panic(0)
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
