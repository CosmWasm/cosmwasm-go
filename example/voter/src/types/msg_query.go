package types

import "github.com/CosmWasm/cosmwasm-go/example/voter/src/state"

// MsgQuery is handled by the Query entrypoint.
type MsgQuery struct {
	// Params returns the current contract parameters.
	Params *struct{} `json:",omitempty"`
	// Voting returns a voting meta.
	Voting *QueryVotingRequest `json:",omitempty"`
	// Results returns a voting summary.
	Tally *QueryTallyRequest `json:",omitempty"`
	// Open returns all open voting IDs.
	Open *struct{} `json:",omitempty"`
	// ReleaseStats returns the current Release operations stats.
	ReleaseStats *struct{} `json:",omitempty"`
	// IBCStats returns sent IBC packets stats for a given senderAddress.
	IBCStats *QueryIBCStatsRequest `json:",omitempty"`
}

// QueryParamsResponse defines MsgQuery.Params response.
type QueryParamsResponse struct {
	state.Params
}

type (
	// QueryVotingRequest defines MsgQuery.Voting request.
	QueryVotingRequest struct {
		ID uint64
	}

	// QueryVotingResponse defines MsgQuery.Voting response.
	QueryVotingResponse struct {
		state.Voting
	}
)

type (
	// QueryTallyRequest defines MsgQuery.Tally request.
	QueryTallyRequest struct {
		ID uint64
	}

	// QueryTallyResponse defines MsgQuery.Tally response.
	QueryTallyResponse struct {
		Open  bool
		Votes []VoteTally
	}

	VoteTally struct {
		Option   string
		TotalYes uint32
		TotalNo  uint32
	}
)

// QueryOpenResponse defines MsgQuery.Open response.
type QueryOpenResponse struct {
	Ids []uint64
}

// QueryReleaseStatsResponse defines MsgQuery.ReleaseStats response.
type QueryReleaseStatsResponse struct {
	state.ReleaseStats
}

type (
	// QueryIBCStatsRequest defines MsgQuery.IBCStats request.
	QueryIBCStatsRequest struct {
		From string
	}

	// QueryIBCStatsResponse defines MsgQuery.IBCStats response.
	QueryIBCStatsResponse struct {
		Stats []state.IBCStats
	}
)
