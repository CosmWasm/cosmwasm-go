package integration

import (
	"testing"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
	mocks "github.com/CosmWasm/wasmvm/api"
	wasmVmTypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ContractTestSuite) TestQueryParams() {
	env := mocks.MockEnv()

	query := types.MsgQuery{Params: &EmptyStruct}
	res, _, err := s.instance.Query(env, query)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	var resp types.QueryParamsResponse
	s.Require().NoError(resp.UnmarshalJSON(res))

	s.Assert().Equal(s.genParams.NewVotingCost.Denom, resp.NewVotingCost.Denom)
	s.Assert().Equal(s.genParams.NewVotingCost.Amount, resp.NewVotingCost.Amount)
	s.Assert().Equal(s.genParams.VoteCost.Denom, resp.VoteCost.Denom)
	s.Assert().Equal(s.genParams.VoteCost.Amount, resp.VoteCost.Amount)
}

func (s *ContractTestSuite) TestQueryVoting() {
	env := mocks.MockEnv()

	s.T().Run("Fail: non-existing", func(t *testing.T) {
		query := types.MsgQuery{
			Voting: &types.QueryVotingRequest{ID: 0},
		}

		_, _, err := s.instance.Query(env, query)
		assert.Error(t, err)
	})

	s.T().Run("OK", func(t *testing.T) {
		// Add voting
		votingID := s.AddVoting(env, s.creatorAddr, "Test", 1000, "a", "b")

		query := types.MsgQuery{
			Voting: &types.QueryVotingRequest{ID: votingID},
		}

		respBz, _, err := s.instance.Query(env, query)
		require.NoError(t, err)

		require.NotNil(t, respBz)
		var resp state.Voting
		require.NoError(t, resp.UnmarshalJSON(respBz))

		assert.Equal(t, votingID, resp.ID)
		assert.Equal(t, s.creatorAddr, resp.CreatorAddr)
		assert.Equal(t, "Test", resp.Name)
		assert.Equal(t, env.Block.Time, resp.StartTime)
		assert.Equal(t, env.Block.Time+1000, resp.EndTime)

		require.Len(t, resp.Tallies, 2)
		assert.Equal(t, "a", resp.Tallies[0].Option)
		assert.Empty(t, resp.Tallies[0].YesAddrs)
		assert.Empty(t, resp.Tallies[0].NoAddrs)
		assert.Equal(t, "b", resp.Tallies[1].Option)
		assert.Empty(t, resp.Tallies[1].YesAddrs)
		assert.Empty(t, resp.Tallies[1].NoAddrs)
	})
}

func (s *ContractTestSuite) TestQueryTally() {
	voter1Addr, voter2Addr := "Voter1Addr", "Voter2Addr"
	env := mocks.MockEnv()

	s.T().Run("Fail: non-existing", func(t *testing.T) {
		query := types.MsgQuery{
			Tally: &types.QueryTallyRequest{ID: 0},
		}

		_, _, err := s.instance.Query(env, query)
		assert.Error(t, err)
	})

	s.T().Run("OK", func(t *testing.T) {
		// Add voting and votes
		votingID := s.AddVoting(env, s.creatorAddr, "Test", 1000, "a", "b")
		env.Block.Time++
		s.Vote(env, voter1Addr, votingID, "a", "yes")
		s.Vote(env, voter2Addr, votingID, "b", "no")

		query := types.MsgQuery{
			Tally: &types.QueryTallyRequest{ID: votingID},
		}

		respBz, _, err := s.instance.Query(env, query)
		require.NoError(t, err)
		require.NotNil(t, respBz)

		var resp types.QueryTallyResponse
		require.NoError(t, resp.UnmarshalJSON(respBz))

		assert.True(t, resp.Open)
		require.Len(t, resp.Votes, 2)

		assert.Equal(t, "a", resp.Votes[0].Option)
		assert.EqualValues(t, 1, resp.Votes[0].TotalYes)
		assert.EqualValues(t, 0, resp.Votes[0].TotalNo)

		assert.Equal(t, "b", resp.Votes[1].Option)
		assert.EqualValues(t, 0, resp.Votes[1].TotalYes)
		assert.EqualValues(t, 1, resp.Votes[1].TotalNo)
	})
}

func (s *ContractTestSuite) TestQueryOpen() {
	env := mocks.MockEnv()
	env.Block.Time = 1

	runQuery := func(t *testing.T) []uint64 {
		query := types.MsgQuery{
			Open: &EmptyStruct,
		}

		respBz, _, err := s.instance.Query(env, query)
		require.NoError(t, err)

		var resp types.QueryOpenResponse
		require.NoError(t, resp.UnmarshalJSON(respBz))

		return resp.Ids
	}

	s.T().Run("No votings", func(t *testing.T) {
		assert.Len(t, runQuery(t), 0)
	})

	// Add votings with different durations
	votingID1 := s.AddVoting(env, s.creatorAddr, "Test1", 10, "a")
	votingID2 := s.AddVoting(env, s.creatorAddr, "Test2", 20, "a")
	votingID3 := s.AddVoting(env, s.creatorAddr, "Test3", 30, "a")

	s.T().Run("3 open", func(t *testing.T) {
		env.Block.Time++

		idsExpected := []uint64{votingID1, votingID2, votingID3}
		idsReceived := runQuery(t)
		assert.ElementsMatch(t, idsExpected, idsReceived)
	})

	s.T().Run("2 open", func(t *testing.T) {
		env.Block.Time = 15

		idsExpected := []uint64{votingID2, votingID3}
		idsReceived := runQuery(t)
		assert.ElementsMatch(t, idsExpected, idsReceived)
	})

	s.T().Run("1 open", func(t *testing.T) {
		env.Block.Time = 25

		idsExpected := []uint64{votingID3}
		idsReceived := runQuery(t)
		assert.ElementsMatch(t, idsExpected, idsReceived)
	})

	s.T().Run("0 open (again)", func(t *testing.T) {
		env.Block.Time = 35

		idsReceived := runQuery(t)
		assert.Empty(t, idsReceived)
	})
}

func (s *ContractTestSuite) TestQueryReleaseStats() {
	env := mocks.MockEnv()
	env.Block.Time = 1

	runQuery := func(t *testing.T) state.ReleaseStats {
		query := types.MsgQuery{
			ReleaseStats: &EmptyStruct,
		}

		respBz, _, err := s.instance.Query(env, query)
		require.NoError(t, err)

		var resp state.ReleaseStats
		require.NoError(t, resp.UnmarshalJSON(respBz))

		return resp
	}

	s.T().Run("No releases", func(t *testing.T) {
		stats := runQuery(t)
		assert.EqualValues(t, 0, stats.Count)
		assert.Nil(t, stats.TotalAmount)
	})

	// Send Release msg and emulate Reply receive
	totalAmtExpected := stdTypes.NewCoinFromUint64(123, "uatom")
	{
		info := mocks.MockInfo(s.creatorAddr, nil)
		releaseMsg := types.MsgExecute{
			Release: &EmptyStruct,
		}

		_, _, err := s.instance.Execute(env, info, releaseMsg)
		s.Require().NoError(err)

		replyMsg := wasmVmTypes.Reply{
			ID: 0,
			Result: wasmVmTypes.SubMsgResult{
				Ok: &wasmVmTypes.SubMsgResponse{
					Events: wasmVmTypes.Events{
						{
							Type: "transfer",
							Attributes: wasmVmTypes.EventAttributes{
								{Key: "amount", Value: totalAmtExpected.String()},
							},
						},
					},
				},
			},
		}

		_, _, err = s.instance.Reply(env, replyMsg)
		s.Require().NoError(err)
	}

	s.T().Run("1 release", func(t *testing.T) {
		stats := runQuery(t)
		assert.EqualValues(t, 1, stats.Count)
		assert.ElementsMatch(t, []stdTypes.Coin{totalAmtExpected}, stats.TotalAmount)
	})
}
