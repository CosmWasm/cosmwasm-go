package integration

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	"github.com/CosmWasm/cosmwasm-go/std/math"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
	"github.com/CosmWasm/cosmwasm-go/systest"
	mocks "github.com/CosmWasm/wasmvm/api"
	wasmVmTypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	EmptyStruct = struct{}{}
)

type ContractTestSuite struct {
	suite.Suite

	instance systest.Instance

	creatorAddr string
	genFunds    []stdTypes.Coin
	genParams   state.Params
}

func (s *ContractTestSuite) SetupTest() {
	t := s.T()

	contractPath := filepath.Join("..", "voter.wasm")
	creatorAddr := "archway1c5qs4le2n6ljgv8fyfy8swu0yf02fh60ccpx75"
	contractFundsCoin := stdTypes.NewCoinFromUint64(1200, "uatom")

	// Load
	instance := systest.NewInstance(t,
		contractPath,
		15_000_000_000_000,
		[]wasmVmTypes.Coin{contractFundsCoin.ToWasmVMCoin()},
	)

	env := mocks.MockEnv()
	info := mocks.MockInfo(creatorAddr, nil)
	msg := types.MsgInstantiate{
		Params: state.Params{
			OwnerAddr: creatorAddr,
			NewVotingCost: stdTypes.Coin{
				Denom:  "uatom",
				Amount: math.NewUint128FromUint64(100),
			},
			VoteCost: stdTypes.Coin{
				Denom:  "uatom",
				Amount: math.NewUint128FromUint64(10),
			},
		},
	}

	// Instantiate
	res, _, err := instance.Instantiate(env, info, msg)
	require.NoError(t, err)

	// Verify response
	require.NotNil(t, res)
	assert.Empty(t, res.Messages)
	assert.Empty(t, res.Attributes)
	assert.Empty(t, res.Events)

	// Setup
	s.instance = instance
	s.creatorAddr = creatorAddr
	s.genFunds = []stdTypes.Coin{contractFundsCoin}
	s.genParams = msg.Params
}

func (s *ContractTestSuite) AddVoting(env wasmVmTypes.Env, creatorAddr, name string, dur uint64, opts ...string) uint64 {
	info := mocks.MockInfo(creatorAddr, []wasmVmTypes.Coin{s.genParams.NewVotingCost.ToWasmVMCoin()})
	msg := types.MsgExecute{
		NewVoting: &types.NewVotingRequest{
			Name:        name,
			VoteOptions: opts,
			Duration:    dur,
		},
	}

	res, _, err := s.instance.Execute(env, info, msg)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.Messages)
	s.Require().Empty(res.Attributes)

	// Verify data
	var resp types.NewVotingResponse
	s.Require().NoError(resp.UnmarshalJSON(res.Data))

	// Verify events
	s.Require().Len(res.Events, 1)
	event := res.Events[0]
	s.Require().Equal(types.EventTypeNewVoting, event.Type)
	s.Require().Len(event.Attributes, 2)

	s.Require().Equal(types.EventAttrKeySender, event.Attributes[0].Key)
	s.Require().Equal(creatorAddr, event.Attributes[0].Value)

	s.Require().Equal(types.EventAttrKeyVotingID, event.Attributes[1].Key)
	s.Require().Equal(strconv.FormatUint(resp.VotingID, 10), event.Attributes[1].Value)

	return resp.VotingID
}

func (s *ContractTestSuite) Vote(env wasmVmTypes.Env, voterAddr string, votingID uint64, opt, vote string) {
	info := mocks.MockInfo(voterAddr, []wasmVmTypes.Coin{s.genParams.VoteCost.ToWasmVMCoin()})
	msg := types.MsgExecute{
		Vote: &types.VoteRequest{
			ID:     votingID,
			Option: opt,
			Vote:   vote,
		},
	}

	res, _, err := s.instance.Execute(env, info, msg)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.Data)
	s.Require().Empty(res.Messages)
	s.Require().Empty(res.Attributes)

	// Verify events
	s.Require().Len(res.Events, 1)
	event := res.Events[0]
	s.Require().Equal(types.EventTypeVote, event.Type)
	s.Require().Len(event.Attributes, 4)

	s.Require().Equal(types.EventAttrKeySender, event.Attributes[0].Key)
	s.Require().Equal(voterAddr, event.Attributes[0].Value)

	s.Require().Equal(types.EventAttrKeyVotingID, event.Attributes[1].Key)
	s.Require().Equal(strconv.FormatUint(votingID, 10), event.Attributes[1].Value)

	s.Require().Equal(types.EventAttrKeyVoteOption, event.Attributes[2].Key)
	s.Require().Equal(opt, event.Attributes[2].Value)

	s.Require().Equal(types.EventAttrKeyVoteDecision, event.Attributes[3].Key)
	s.Require().Equal(vote, event.Attributes[3].Value)
}

func TestContract(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}
