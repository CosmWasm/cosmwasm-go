package src

import (
	"encoding/json"
	"testing"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/stretchr/testify/require"
)

// This demos a simple call of one function with a mock environment which can be queried after the call
func TestInit(t *testing.T) {
	cases := map[string]struct {
		initMsg []byte
		funds   []std.Coin
		valid   bool
	}{
		// TODO: why doesn't ezjson.Unmarshal return an error here?
		"invalid json": {initMsg: []byte("{...")},
		"wrong struct": {initMsg: []byte(`{"foo": 1, "bar": "world"}`)},
		"proper init":  {initMsg: []byte(`{"name":"Cool Coin","symbol":"COOL","decimal":6,"total_supply":12345678}`), valid: true},
		"proper init with funds": {
			initMsg: []byte(`{"name":"Cool Coin","symbol":"COOL","decimal":6,"total_supply":12345678}`),
			valid:   true,
			funds:   []std.Coin{{Denom: "uatom", Amount: "1000"}},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			deps := std.MockDeps()
			env := std.MockEnv()
			info := std.MockInfo("creator", tc.funds)
			res, err := Init(deps, env, info, tc.initMsg)
			if tc.valid {
				require.Nil(t, err)
				require.NotNil(t, res)
				require.Equal(t, 1, len(res.Ok.Attributes))
				require.Equal(t, "hello", res.Ok.Attributes[0].Key)
				require.Equal(t, "world", res.Ok.Attributes[0].Value)
				// make sure we wrote the owner
				owner := NewOwnership(deps)
				owner.LoadOwner()
				require.NotEmpty(t, owner.GetOwner())
				require.Empty(t, owner.GetNewOwner())
			} else {
				require.NotNil(t, err)
				require.Nil(t, res)
			}
		})
	}
}

// This demos unit testing state objects
func TestOwner(t *testing.T) {
	deps := std.MockDeps()
	ownership := NewOwnership(deps)

	// write some data
	first := []byte("original")
	second := []byte("updated")
	other := []byte("random")

	// set owner
	ownership.Owned(first)
	require.Equal(t, first, ownership.GetOwner())
	require.Empty(t, ownership.GetNewOwner())

	// transfer owner fails
	ownership.TransferOwnership(other, second)
	require.Equal(t, first, ownership.GetOwner())
	require.Empty(t, ownership.GetNewOwner())

	// transfer owner succeeds
	ownership.TransferOwnership(first, second)
	require.Equal(t, first, ownership.GetOwner())
	require.Equal(t, second, ownership.GetNewOwner())

	// only the new owner can accept
	ownership.AcceptTransfer(other, second)
	require.Equal(t, first, ownership.GetOwner())
	require.Equal(t, second, ownership.GetNewOwner())

	// this works
	ownership.AcceptTransfer(second, second)
	require.Equal(t, second, ownership.GetOwner())
	require.Empty(t, ownership.GetNewOwner())

	// set some other new owner to ensure save/load work
	ownership.TransferOwnership(second, first)
	require.Equal(t, second, ownership.GetOwner())
	require.Equal(t, first, ownership.GetNewOwner())

	// another object cannot see this state yet
	loaded := NewOwnership(deps)
	loaded.LoadOwner()
	require.Empty(t, loaded.GetOwner())
	require.Empty(t, loaded.GetNewOwner())

	// save and load
	ownership.SaveOwner()
	loaded.LoadOwner()
	// data is available
	require.Equal(t, second, loaded.GetOwner())
	require.Equal(t, first, loaded.GetNewOwner())
}

// This runs the same path we do in cosmwasm-simulate
func TestWorkflow(t *testing.T) {
	deps := std.MockDeps()
	env := std.MockEnv()
	info := std.MockInfo("original_owner_addr", nil)

	initMsg := []byte(`{"name":"OKB","symbol":"OKB","decimal":10,"total_supply":170000}`)
	ires, err := Init(deps, env, info, initMsg)
	require.Nil(t, err)
	require.NotNil(t, ires)

	handleMsg := []byte(`{"Transfer":{"to":"1234567","value": 2000}}`)
	hres, err := Invoke(deps, env, info, handleMsg)
	require.Nil(t, err)
	require.NotNil(t, hres)

	queryMsg := []byte(`{"balance":{"address":"1234567"}}`)
	qres, err := Query(deps, env, queryMsg)
	require.Nil(t, err)
	require.NotEmpty(t, qres.Ok)

	// let us parse the query??
	var bal BalanceResponse
	jerr := json.Unmarshal(qres.Ok, &bal)
	require.NoError(t, jerr)
	require.Equal(t, uint64(2000), bal.Value)
}
