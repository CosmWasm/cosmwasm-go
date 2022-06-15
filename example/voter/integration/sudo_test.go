package integration

import (
	"testing"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
	mocks "github.com/CosmWasm/wasmvm/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ContractTestSuite) TestSudoChangeAddVotingCost() {
	env := mocks.MockEnv()
	expectedCoin := stdTypes.Coin{
		Denom:  s.genParams.NewVotingCost.Denom,
		Amount: s.genParams.NewVotingCost.Amount.Sub64(1),
	}

	s.T().Run("OK", func(t *testing.T) {
		msg := types.MsgSudo{
			ChangeNewVotingCost: &types.ChangeCostRequest{
				NewCost: expectedCoin,
			},
		}

		res, _, err := s.instance.Sudo(env, msg)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Empty(t, res.Data)
		assert.Empty(t, res.Messages)
		assert.Empty(t, res.Attributes)

		// Verify events
		require.Len(t, res.Events, 1)
		rcvEvent := res.Events[0]
		assert.Equal(t, types.EventTypeNewVotingCostChanged, rcvEvent.Type)

		require.Len(t, rcvEvent.Attributes, 2)
		assert.Equal(t, types.EventAttrKeyOldCost, rcvEvent.Attributes[0].Key)
		assert.Equal(t, s.genParams.NewVotingCost.String(), rcvEvent.Attributes[0].Value)
		assert.Equal(t, types.EventAttrKeyNewCost, rcvEvent.Attributes[1].Key)
		assert.Equal(t, expectedCoin.String(), rcvEvent.Attributes[1].Value)

		// Verify state change
		query := types.MsgQuery{Params: &EmptyStruct}
		paramsBz, _, err := s.instance.Query(env, query)
		require.NoError(t, err)
		require.NotNil(t, paramsBz)

		var params types.QueryParamsResponse
		require.NoError(t, params.UnmarshalJSON(paramsBz))
		assert.Equal(t, expectedCoin, params.NewVotingCost)
	})

	s.T().Run("Fail: invalid input", func(t *testing.T) {
		msg := types.MsgSudo{
			ChangeNewVotingCost: &types.ChangeCostRequest{
				NewCost: stdTypes.Coin{
					Denom:  "1uatom",
					Amount: expectedCoin.Amount,
				},
			},
		}

		_, _, err := s.instance.Sudo(env, msg)
		require.Error(t, err)
	})
}

func (s *ContractTestSuite) TestSudoChangeVoteCost() {
	env := mocks.MockEnv()
	expectedCoin := stdTypes.Coin{
		Denom:  s.genParams.VoteCost.Denom,
		Amount: s.genParams.VoteCost.Amount.Sub64(1),
	}

	s.T().Run("OK", func(t *testing.T) {
		msg := types.MsgSudo{
			ChangeVoteCost: &types.ChangeCostRequest{
				NewCost: expectedCoin,
			},
		}

		res, _, err := s.instance.Sudo(env, msg)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Empty(t, res.Data)
		assert.Empty(t, res.Messages)
		assert.Empty(t, res.Attributes)

		// Verify events
		require.Len(t, res.Events, 1)
		rcvEvent := res.Events[0]
		assert.Equal(t, types.EventTypeVoteCostChanged, rcvEvent.Type)

		require.Len(t, rcvEvent.Attributes, 2)
		assert.Equal(t, types.EventAttrKeyOldCost, rcvEvent.Attributes[0].Key)
		assert.Equal(t, s.genParams.VoteCost.String(), rcvEvent.Attributes[0].Value)
		assert.Equal(t, types.EventAttrKeyNewCost, rcvEvent.Attributes[1].Key)
		assert.Equal(t, expectedCoin.String(), rcvEvent.Attributes[1].Value)

		// Verify state change
		query := types.MsgQuery{Params: &EmptyStruct}
		paramsBz, _, err := s.instance.Query(env, query)
		require.NoError(t, err)
		require.NotNil(t, paramsBz)

		var params types.QueryParamsResponse
		require.NoError(t, params.UnmarshalJSON(paramsBz))
		assert.Equal(t, expectedCoin, params.VoteCost)
	})

	s.T().Run("Fail: invalid input", func(t *testing.T) {
		msg := types.MsgSudo{
			ChangeVoteCost: &types.ChangeCostRequest{
				NewCost: stdTypes.Coin{
					Denom:  "1uatom",
					Amount: expectedCoin.Amount,
				},
			},
		}

		_, _, err := s.instance.Sudo(env, msg)
		require.Error(t, err)
	})
}
