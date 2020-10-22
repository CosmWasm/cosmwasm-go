package src

import (
	"testing"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	cases := map[string]struct {
		initMsg []byte
		valid   bool
	}{
		"invalid json": {initMsg: []byte("{...")},
		"wrong struct": {initMsg: []byte(`{"foo": 1, "bar": "world"}`)},
		"proper init":  {initMsg: []byte(`{"name":"Cool Coin","symbol":"COOL","decimal":6,"total_supply":12345678}`)},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			deps := std.MockExtern()
			env := std.MockEnv("creator", nil)
			res, err := Init(deps, env, tc.initMsg)
			if tc.valid {
				require.Nil(t, err)
				require.NotNil(t, res)
			} else {
				require.NotNil(t, err)
				require.Nil(t, res)
			}
		})
	}
}
