package src

import (
	"encoding/json"
	"testing"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustEncode(t *testing.T, msg interface{}) []byte {
	bz, err := json.Marshal(msg)
	require.NoError(t, err)
	return bz
}

const VERIFIER = "VERIFIER"
const BENEFICIARY = "benefits"
const FUNDER = "creator"

// this can be used for a quick setup if you don't have nay other requirements
func defaultInit(t *testing.T) *std.Deps {
	deps := std.MockDeps()
	env := std.MockEnv()
	info := std.MockInfo(FUNDER, nil)
	initMsg := InitMsg{
		Verifier:    VERIFIER,
		Beneficiary: BENEFICIARY,
	}
	res, err := Init(deps, env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	return deps
}

func TestInitAndQuery(t *testing.T) {
	deps := std.MockDeps()
	env := std.MockEnv()
	info := std.MockInfo(FUNDER, nil)
	initMsg := InitMsg{
		Verifier:    VERIFIER,
		Beneficiary: BENEFICIARY,
	}
	res, err := Init(deps, env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 0, len(res.Ok.Messages))
	assert.Equal(t, 1, len(res.Ok.Attributes))
	attr := res.Ok.Attributes[0]
	assert.Equal(t, "Let the", attr.Key)
	assert.Equal(t, "hacking begin", attr.Value)

	qmsg := []byte(`{"VERIFIER":{}}`)
	data, err := Query(deps, env, qmsg)
	require.NoError(t, err)
	var qres VerifierResponse
	err = json.Unmarshal(data.Ok, &qres)
	require.NoError(t, err)
	assert.Equal(t, VERIFIER, qres.Verifier)
}

func TestPanic(t *testing.T) {
	deps := defaultInit(t)
	env := std.MockEnv()
	info := std.MockInfo(FUNDER, nil)
	handleMsg := []byte(`{"panic":{}}`)
	require.Panics(t, func() {
		Handle(deps, env, info, handleMsg)
	})
}

func TestRelease(t *testing.T) {
	cases := map[string]struct {
		signer string
		valid  bool
	}{
		"verifier releases": {VERIFIER, true},
		"random fails":      {BENEFICIARY, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// TODO: figure out how to set query value and then query from the contract
			deps := defaultInit(t)
			env := std.MockEnv()
			info := std.MockInfo(tc.signer, nil)
			handleMsg := []byte(`{"release":{}}`)
			res, err := Handle(deps, env, info, handleMsg)
			if !tc.valid {
				require.Error(t, err)
				require.Equal(t, "Unauthorized", err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Equal(t, 1, len(res.Ok.Messages))
				msg := res.Ok.Messages[0]
				expected := std.CosmosMsg{Bank: std.BankMsg{Send: std.SendMsg{
					FromAddress: std.MOCK_CONTRACT_ADDR,
					ToAddress:   BENEFICIARY,
					// TODO: get this from the validator
					Amount: []std.Coin{{"ATOM", "1000"}},
				}}}
				assert.Equal(t, expected, msg)
				assert.Equal(t, 2, len(res.Ok.Attributes))
				assert.Equal(t, []std.Attribute{{"action", "release"}, {"destination", BENEFICIARY}}, res.Ok.Attributes)
			}
		})
	}
}

//
//func TestHandle(t *testing.T) {
//	deps := std.MockDeps()
//	env := std.MockEnv()
//	info := std.MockInfo("creator", nil)
//	initMsg := []byte(`{"count":123}`)
//	res, err := Init(deps, env, info, initMsg)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	info = std.MockInfo("random", nil)
//	handleMsg := []byte(`{"increment":{}}`)
//	_, err = Handle(deps, env, info, handleMsg)
//	require.NoError(t, err)
//
//	qmsg := []byte(`{"get_count":{}}`)
//	data, err := Query(deps, env, qmsg)
//	require.NoError(t, err)
//	var qres CountResponse
//	err = json.Unmarshal(data.Ok, &qres)
//	require.NoError(t, err)
//	require.Equal(t, uint64(124), qres.Count)
//}
