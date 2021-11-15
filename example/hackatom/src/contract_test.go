package src

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/mocks"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

func mustEncode(t *testing.T, msg interface{}) []byte {
	bz, err := json.Marshal(msg)
	require.NoError(t, err)
	return bz
}

const VERIFIER = "verifies"
const BENEFICIARY = "benefits"
const FUNDER = "creator"

// this can be used for a quick setup if you don't have nay other requirements
func defaultInit(t *testing.T, funds []types.Coin) *std.Deps {
	deps := mocks.MockDeps(funds)
	env := mocks.MockEnv()
	info := mocks.MockInfo(FUNDER, funds)
	initMsg := InitMsg{
		Verifier:    VERIFIER,
		Beneficiary: BENEFICIARY,
	}
	res, err := Instantiate(deps, env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	return deps
}

func TestInitAndQuery(t *testing.T) {
	deps := mocks.MockDeps(nil)
	env := mocks.MockEnv()
	info := mocks.MockInfo(FUNDER, nil)
	initMsg := InitMsg{
		Verifier:    VERIFIER,
		Beneficiary: BENEFICIARY,
	}
	res, err := Instantiate(deps, env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 0, len(res.Messages))
	require.Equal(t, 1, len(res.Attributes))
	attr := res.Attributes[0]
	assert.Equal(t, "Let the", attr.Key)
	assert.Equal(t, "hacking begin", attr.Value)

	qmsg := []byte(`{"verifier":{}}`)
	data, err := Query(deps, env, qmsg)
	require.NoError(t, err)
	var qres VerifierResponse
	require.NoError(t, err)
	err = json.Unmarshal(data, &qres)
	require.NoError(t, err)
	assert.Equal(t, VERIFIER, qres.Verifier)
}

func TestPanic(t *testing.T) {
	deps := defaultInit(t, nil)
	env := mocks.MockEnv()
	info := mocks.MockInfo(FUNDER, nil)
	handleMsg := []byte(`{"panic":{}}`)
	require.Panics(t, func() {
		_, _ = Execute(deps, env, info, handleMsg)
	})
}

func TestRelease(t *testing.T) {
	cases := map[string]struct {
		signer string
		funds  []types.Coin
		valid  bool
	}{
		"verifier releases": {VERIFIER, types.NewCoins(765432, "wei"), true},
		"random fails":      {BENEFICIARY, types.NewCoins(765432, "wei"), false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// TODO: figure out how to set query value and then query from the contract
			deps := defaultInit(t, tc.funds)
			env := mocks.MockEnv()
			info := mocks.MockInfo(tc.signer, nil)
			handleMsg := []byte(`{"release":{}}`)
			res, err := Execute(deps, env, info, handleMsg)
			if !tc.valid {
				require.Error(t, err)
				require.Equal(t, "Unauthorized", err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Equal(t, 1, len(res.Messages))
				msg := res.Messages[0]
				expected := types.CosmosMsg{Bank: &types.BankMsg{Send: &types.SendMsg{
					ToAddress: BENEFICIARY,
					Amount:    tc.funds,
				}}}
				assert.Equal(t, expected, msg.Msg)
				assert.Equal(t, 2, len(res.Attributes))
				assert.Equal(t, []types.EventAttribute{{"action", "release"}, {"destination", BENEFICIARY}}, res.Attributes)
			}
		})
	}
}
