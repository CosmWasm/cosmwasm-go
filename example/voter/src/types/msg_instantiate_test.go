package types

import (
	"testing"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgInstantiateValidate(t *testing.T) {
	type testCase struct {
		name string
		info stdTypes.MessageInfo
		msg  MsgInstantiate
		//
		errExpected bool
	}

	senderAccAddr := "SenderAccAddress"
	otherAccAddr := "OtherAccAddress"

	testCases := []testCase{
		{
			name: "OK",
			info: stdTypes.MessageInfo{Sender: senderAccAddr},
			msg: MsgInstantiate{
				Params: state.Params{
					OwnerAddr:     senderAccAddr,
					NewVotingCost: stdTypes.NewCoinFromUint64(100, "uatom"),
					VoteCost:      stdTypes.NewCoinFromUint64(100, "uatom"),
				},
			},
		},
		{
			name: "Fail: OwnerAddr: mismatch",
			info: stdTypes.MessageInfo{Sender: otherAccAddr},
			msg: MsgInstantiate{
				Params: state.Params{
					OwnerAddr:     senderAccAddr,
					NewVotingCost: stdTypes.NewCoinFromUint64(100, "uatom"),
					VoteCost:      stdTypes.NewCoinFromUint64(100, "uatom"),
				},
			},
			errExpected: true,
		},
		{
			name: "Fail: NewVotingCost: invalid denom",
			info: stdTypes.MessageInfo{Sender: senderAccAddr},
			msg: MsgInstantiate{
				Params: state.Params{
					OwnerAddr:     senderAccAddr,
					NewVotingCost: stdTypes.NewCoinFromUint64(100, "1uatom"),
					VoteCost:      stdTypes.NewCoinFromUint64(100, "uatom"),
				},
			},
			errExpected: true,
		},
		{
			name: "Fail: VoteCost: invalid denom",
			info: stdTypes.MessageInfo{Sender: senderAccAddr},
			msg: MsgInstantiate{
				Params: state.Params{
					OwnerAddr:     senderAccAddr,
					NewVotingCost: stdTypes.NewCoinFromUint64(100, "uatom"),
					VoteCost:      stdTypes.NewCoinFromUint64(100, "1uatom"),
				},
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errExpected {
				assert.Error(t, tc.msg.Validate(tc.info))
				return
			}
			assert.NoError(t, tc.msg.Validate(tc.info))
		})
	}
}
