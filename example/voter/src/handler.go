package src

import (
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/pkg"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	"github.com/CosmWasm/cosmwasm-go/std"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
)

// handleMsgInstantiate handles types.MsgInstantiate msg.
func handleMsgInstantiate(deps *std.Deps, info stdTypes.MessageInfo, msg types.MsgInstantiate) (*stdTypes.Response, error) {
	// Input check
	if err := msg.Validate(info); err != nil {
		return nil, types.NewErrInvalidRequest("msg validation: " + err.Error())
	}

	// Set initial contract state
	if err := state.SetParams(deps.Storage, msg.Params); err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	return &stdTypes.Response{}, nil
}

// handleMsgRelease handles MsgExecute.Release msg.
func handleMsgRelease(deps *std.Deps, env stdTypes.Env, info stdTypes.MessageInfo) (*stdTypes.Response, error) {
	// Input check
	params, err := state.GetParams(deps.Storage)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	if info.Sender != params.OwnerAddr {
		return nil, types.NewErrInvalidRequest("release can be done only by the contract owner")
	}

	// Transfer
	queryClient := std.QuerierWrapper{Querier: deps.Querier}
	contractFunds, err := queryClient.QueryAllBalances(env.Contract.Address)
	if err != nil {
		return nil, types.NewErrInternal("bank balance query: " + err.Error())
	}

	bankMsg := stdTypes.NewSubMsg(stdTypes.SendMsg{
		ToAddress: info.Sender,
		Amount:    contractFunds,
	})

	// Result
	res := types.ReleaseResponse{
		ReleasedAmount: contractFunds,
	}

	resBz, err := res.MarshalJSON()
	if err != nil {
		return nil, types.NewErrInternal("result JSON marshal: " + err.Error())
	}

	return &stdTypes.Response{
		Data: resBz,
		Messages: []stdTypes.SubMsg{
			bankMsg,
		},
		Events: []stdTypes.Event{
			types.NewEventRelease(info.Sender),
		},
	}, nil
}

// handleMsgNewVoting handles MsgExecute.NewVoting msg.
func handleMsgNewVoting(deps *std.Deps, env stdTypes.Env, info stdTypes.MessageInfo, req types.NewVotingRequest) (*stdTypes.Response, error) {
	// Input check
	if err := req.Validate(); err != nil {
		return nil, types.NewErrInvalidRequest("req validation: " + err.Error())
	}

	params, err := state.GetParams(deps.Storage)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	if err := pkg.CoinsContainMinAmount(info.Funds, params.NewVotingCost); err != nil {
		return nil, types.NewErrInvalidRequest(err.Error())
	}

	// Create a new voting
	votingID, err := state.NextVotingID(deps.Storage)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	voting := state.NewVoting(votingID, req.Name, info.Sender, env.Block.Time, req.Duration, req.VoteOptions)

	// Update contract state
	state.SetLastVotingID(deps.Storage, votingID)
	if err := state.SetVoting(deps.Storage, voting); err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	// Result
	res := types.NewVotingResponse{
		VotingID: votingID,
	}

	resBz, err := res.MarshalJSON()
	if err != nil {
		return nil, types.NewErrInternal("result JSON marshal: " + err.Error())
	}

	return &stdTypes.Response{
		Data: resBz,
		Events: []stdTypes.Event{
			types.NewEventVotingCreated(info.Sender, votingID),
		},
	}, nil
}

// handleMsgVote handles MsgExecute.Vote msg.
func handleMsgVote(deps *std.Deps, env stdTypes.Env, info stdTypes.MessageInfo, req types.VoteRequest) (*stdTypes.Response, error) {
	// Input check
	if err := req.Validate(); err != nil {
		return nil, types.NewErrInvalidRequest("req validation: " + err.Error())
	}

	params, err := state.GetParams(deps.Storage)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	if err := pkg.CoinsContainMinAmount(info.Funds, params.VoteCost); err != nil {
		return nil, types.NewErrInvalidRequest(err.Error())
	}

	voting, err := state.GetVoting(deps.Storage, req.ID)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}
	if voting == nil {
		return nil, types.NewErrInvalidRequest("voting with requested ID not found")
	}

	if voting.IsClosed(env.Block.Time) {
		return nil, types.ErrVotingClosed
	}
	if voting.HasVote(info.Sender) {
		return nil, types.ErrAlreadyVoted
	}

	// Append vote and update contract state
	var voteErr error
	switch req.Vote {
	case "yes":
		voteErr = voting.AddYesVote(req.Option, info.Sender)
	case "no":
		voteErr = voting.AddNoVote(req.Option, info.Sender)
	}
	if voteErr != nil {
		return nil, types.NewErrInvalidRequest(voteErr.Error())
	}

	if err := state.SetVoting(deps.Storage, *voting); err != nil {
		return nil, types.NewErrInternal(voteErr.Error())
	}

	return &stdTypes.Response{
		Events: []stdTypes.Event{
			types.NewEventVoteAdded(info.Sender, req.ID, req.Option, req.Vote),
		},
	}, nil
}

// handleSudoChangeNewVotingCost handles MsgSudo.NewVotingCost msg.
func handleSudoChangeNewVotingCost(deps *std.Deps, req types.ChangeCostRequest) (*stdTypes.Response, error) {
	// Input check
	if err := req.Validate(); err != nil {
		return nil, types.NewErrInvalidRequest("req validation: " + err.Error())
	}

	// Update params state
	params, err := state.GetParams(deps.Storage)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}
	oldCost := params.NewVotingCost

	params.NewVotingCost = req.NewCost
	if err := state.SetParams(deps.Storage, params); err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	return &stdTypes.Response{
		Events: []stdTypes.Event{
			types.NewEventNewVotingCostChanged(oldCost, req.NewCost),
		},
	}, nil
}

// handleSudoChangeVoteCost handles MsgSudo.VoteCost msg.
func handleSudoChangeVoteCost(deps *std.Deps, req types.ChangeCostRequest) (*stdTypes.Response, error) {
	// Input check
	if err := req.Validate(); err != nil {
		return nil, types.NewErrInvalidRequest("req validation: " + err.Error())
	}

	// Update params state
	params, err := state.GetParams(deps.Storage)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}
	oldCost := params.VoteCost

	params.VoteCost = req.NewCost
	if err := state.SetParams(deps.Storage, params); err != nil {
		return nil, types.NewErrInternal(err.Error())
	}

	return &stdTypes.Response{
		Events: []stdTypes.Event{
			types.NewEventVoteCostChanged(oldCost, req.NewCost),
		},
	}, nil
}
