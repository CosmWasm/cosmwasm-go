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

func (s *ContractTestSuite) TestReplyBankSend() {
	env := mocks.MockEnv()

	// Add voting
	s.AddVoting(env, s.creatorAddr, "Test", 100, "a")

	s.T().Run("Fail: no reply ID found", func(t *testing.T) {
		msg := wasmVmTypes.Reply{
			ID: 0,
		}

		_, _, err := s.instance.Reply(env, msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	// Release funds (replyID 0 is created here)
	{
		info := mocks.MockInfo(s.creatorAddr, nil)
		msg := types.MsgExecute{
			Release: &EmptyStruct,
		}

		_, _, err := s.instance.Execute(env, info, msg)
		s.Require().NoError(err)
	}

	s.T().Run("Fail: invalid reply: with error", func(t *testing.T) {
		msg := wasmVmTypes.Reply{
			ID: 0,
			Result: wasmVmTypes.SubMsgResult{
				Err: "some error",
			},
		}

		_, _, err := s.instance.Reply(env, msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "x/bank reply: error received")
	})

	s.T().Run("Fail: invalid reply: invalid message type received (wrong events)", func(t *testing.T) {
		msg := wasmVmTypes.Reply{
			ID: 0,
			Result: wasmVmTypes.SubMsgResult{
				Ok: &wasmVmTypes.SubMsgResponse{},
			},
		}

		_, _, err := s.instance.Reply(env, msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "x/bank reply: transfer.amount attribute: not found")
	})

	s.T().Run("OK", func(t *testing.T) {
		releaseAmtExpected := stdTypes.NewCoinFromUint64(1000, "uatom")

		msg := wasmVmTypes.Reply{
			ID: 0,
			Result: wasmVmTypes.SubMsgResult{
				Ok: &wasmVmTypes.SubMsgResponse{
					Events: wasmVmTypes.Events{
						{
							Type: "transfer",
							Attributes: wasmVmTypes.EventAttributes{
								{Key: "amount", Value: releaseAmtExpected.String()},
							},
						},
					},
				},
			},
		}

		_, _, err := s.instance.Reply(env, msg)
		assert.NoError(t, err)

		// Verify stats changed
		{
			query := types.MsgQuery{
				ReleaseStats: &EmptyStruct,
			}

			respBz, _, err := s.instance.Query(env, query)
			require.NoError(t, err)
			require.NotNil(t, respBz)

			var stats state.ReleaseStats
			require.NoError(t, stats.UnmarshalJSON(respBz))

			assert.EqualValues(t, 1, stats.Count)
			assert.ElementsMatch(t, []stdTypes.Coin{releaseAmtExpected}, stats.TotalAmount)
		}
	})
}
