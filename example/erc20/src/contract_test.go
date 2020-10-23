package src

import (
	"testing"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/stretchr/testify/require"
)

// This demos a simple call of one function with a mock environment which can be queried after the call
func TestInit(t *testing.T) {
	cases := map[string]struct {
		initMsg []byte
		valid   bool
	}{
		// TODO: why doesn't exjson.Unmarshal return an error here?
		"invalid json": {initMsg: []byte("{...")},
		"wrong struct": {initMsg: []byte(`{"foo": 1, "bar": "world"}`)},
		"proper init":  {initMsg: []byte(`{"name":"Cool Coin","symbol":"COOL","decimal":6,"total_supply":12345678}`), valid: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			deps := std.MockExtern()
			env := std.MockEnv("creator", nil)
			res, err := Init(deps, env, tc.initMsg)
			if tc.valid {
				require.Nil(t, err)
				require.NotNil(t, res)
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
	deps := std.MockExtern()
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
