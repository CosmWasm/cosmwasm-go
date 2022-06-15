package src

import (
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	"github.com/CosmWasm/cosmwasm-go/std"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
)

// queryParams handles MsgQuery.Params query.
func queryParams(deps *std.Deps) (*types.QueryParamsResponse, error) {
	params, err := state.GetParams(deps.Storage)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// queryVoting handles MsgQuery.Voting query.
func queryVoting(deps *std.Deps, req types.QueryVotingRequest) (*types.QueryVotingResponse, error) {
	voting, err := state.GetVoting(deps.Storage, req.ID)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	if voting == nil {
		return nil, stdTypes.NotFound{Kind: "voting"}
	}

	return &types.QueryVotingResponse{
		Voting: *voting,
	}, nil
}

// queryTally handles MsgQuery.Tally query.
func queryTally(deps *std.Deps, env stdTypes.Env, req types.QueryTallyRequest) (*types.QueryTallyResponse, error) {
	voting, err := state.GetVoting(deps.Storage, req.ID)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	if voting == nil {
		return nil, stdTypes.NotFound{Kind: "voting"}
	}

	resp := types.QueryTallyResponse{
		Open:  !voting.IsClosed(env.Block.Time),
		Votes: make([]types.VoteTally, 0, len(voting.Tallies)),
	}

	for _, tally := range voting.Tallies {
		resp.Votes = append(
			resp.Votes,
			types.VoteTally{
				Option:   tally.Option,
				TotalYes: uint32(len(tally.YesAddrs)),
				TotalNo:  uint32(len(tally.NoAddrs)),
			},
		)
	}

	return &resp, nil
}

// queryOpen handles MsgQuery.Open query.
func queryOpen(deps *std.Deps, env stdTypes.Env) (*types.QueryOpenResponse, error) {
	ids := make([]uint64, 0)

	err := state.IterateVotings(deps.Storage, func(voting state.Voting) (stop bool) {
		if !voting.IsClosed(env.Block.Time) {
			ids = append(ids, voting.ID)
		}
		return false
	})
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	return &types.QueryOpenResponse{
		Ids: ids,
	}, nil
}
