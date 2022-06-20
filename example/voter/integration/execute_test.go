package integration

import (
	"testing"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	mocks "github.com/CosmWasm/wasmvm/api"
	wasmVmTypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ContractTestSuite) TestExecuteNewVoting() {
	env := mocks.MockEnv()

	// Test OK
	s.AddVoting(env, s.creatorAddr, "Test", 100, "a")

	s.T().Run("Fail: invalid input", func(t *testing.T) {
		info := mocks.MockInfo(s.creatorAddr, []wasmVmTypes.Coin{s.genParams.NewVotingCost.ToWasmVMCoin()})
		msg := types.MsgExecute{
			NewVoting: &types.NewVotingRequest{
				Name:        "Test",
				VoteOptions: nil,
				Duration:    100,
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("Fail: invalid payment", func(t *testing.T) {
		payment := s.genParams.NewVotingCost
		payment.Amount = payment.Amount.Sub64(1)

		info := mocks.MockInfo(s.creatorAddr, []wasmVmTypes.Coin{payment.ToWasmVMCoin()})
		msg := types.MsgExecute{
			NewVoting: &types.NewVotingRequest{
				Name:        "Test",
				VoteOptions: []string{"a"},
				Duration:    100,
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})
}

func (s *ContractTestSuite) TestExecuteVote() {
	env := mocks.MockEnv()

	voter1Addr, voter2Addr := "Voter1Addr", "Voter2Addr"

	// Test OK
	votingID := s.AddVoting(env, s.creatorAddr, "Test", 100, "a")
	s.Vote(env, voter1Addr, votingID, "a", "yes")

	s.T().Run("Fail: invalid input", func(t *testing.T) {
		info := mocks.MockInfo(voter2Addr, []wasmVmTypes.Coin{s.genParams.VoteCost.ToWasmVMCoin()})
		msg := types.MsgExecute{
			Vote: &types.VoteRequest{
				ID:     votingID,
				Option: "",
				Vote:   "yes",
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("Fail: invalid payment", func(t *testing.T) {
		payment := s.genParams.VoteCost
		payment.Amount = payment.Amount.Sub64(1)

		info := mocks.MockInfo(voter2Addr, []wasmVmTypes.Coin{payment.ToWasmVMCoin()})
		msg := types.MsgExecute{
			Vote: &types.VoteRequest{
				ID:     votingID,
				Option: "a",
				Vote:   "yes",
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("Fail: non-existing voting", func(t *testing.T) {
		info := mocks.MockInfo(voter2Addr, []wasmVmTypes.Coin{s.genParams.VoteCost.ToWasmVMCoin()})
		msg := types.MsgExecute{
			Vote: &types.VoteRequest{
				ID:     votingID + 1,
				Option: "a",
				Vote:   "yes",
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("Fail: already voted", func(t *testing.T) {
		info := mocks.MockInfo(voter1Addr, []wasmVmTypes.Coin{s.genParams.VoteCost.ToWasmVMCoin()})
		msg := types.MsgExecute{
			Vote: &types.VoteRequest{
				ID:     votingID,
				Option: "a",
				Vote:   "no",
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("Fail: voting is closed", func(t *testing.T) {
		env := mocks.MockEnv()
		env.Block.Time += 200
		info := mocks.MockInfo(voter2Addr, []wasmVmTypes.Coin{s.genParams.VoteCost.ToWasmVMCoin()})
		msg := types.MsgExecute{
			Vote: &types.VoteRequest{
				ID:     votingID,
				Option: "a",
				Vote:   "yes",
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("Fail: non-existing option", func(t *testing.T) {
		info := mocks.MockInfo(voter2Addr, []wasmVmTypes.Coin{s.genParams.VoteCost.ToWasmVMCoin()})
		msg := types.MsgExecute{
			Vote: &types.VoteRequest{
				ID:     votingID,
				Option: "c",
				Vote:   "no",
			},
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})
}

func (s *ContractTestSuite) TestExecuteRelease() {
	env := mocks.MockEnv()

	voter1Addr, voter2Addr := "Voter1Addr", "Voter2Addr"

	// Add voting and votes (1000 + 2 * 100 of raised funds)
	votingID := s.AddVoting(env, voter1Addr, "Test", 100, "a")
	s.Vote(env, voter1Addr, votingID, "a", "yes")
	s.Vote(env, voter2Addr, votingID, "a", "no")

	s.T().Run("Fail: unauthorized", func(t *testing.T) {
		info := mocks.MockInfo(voter2Addr, nil)
		msg := types.MsgExecute{
			Release: &EmptyStruct,
		}

		_, _, err := s.instance.Execute(env, info, msg)
		assert.Error(t, err)
	})

	s.T().Run("OK", func(t *testing.T) {
		info := mocks.MockInfo(s.creatorAddr, nil)
		msg := types.MsgExecute{
			Release: &EmptyStruct,
		}

		res, _, err := s.instance.Execute(env, info, msg)
		require.NoError(t, err)
		require.NotNil(t, res)

		var resp types.ReleaseResponse
		require.NoError(t, resp.UnmarshalJSON(res.Data))

		assert.ElementsMatch(t, s.genFunds, resp.ReleasedAmount)
	})
}
