package src

import (
	"encoding/json"
	"testing"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	cases := map[string]struct {
		initMsg []byte
		count   uint64
	}{
		// TODO: why doesn't ezjson.Unmarshal return an error here?
		"random json": {initMsg: []byte("{...")},
		"count zero":  {initMsg: []byte(`{"count":0}`)},
		"count large": {initMsg: []byte(`{"count":123}`), count: 123},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			deps := std.MockDeps()
			env := std.MockEnv()
			info := std.MockInfo("creator", nil)
			res, err := Init(deps, env, info, tc.initMsg)
			require.NoError(t, err)
			require.NotNil(t, res)
			// check any return value if needed

			qmsg := []byte(`{"get_count":{}}`)
			data, err := Query(deps, env, qmsg)
			require.NoError(t, err)
			var qres CountResponse
			err = json.Unmarshal(data.Ok, &qres)
			require.NoError(t, err)
			require.Equal(t, tc.count, qres.Count)
		})
	}

}

func TestHandle(t *testing.T) {
	deps := std.MockDeps()
	env := std.MockEnv()
	info := std.MockInfo("creator", nil)
	initMsg := []byte(`{"count":123}`)
	res, err := Init(deps, env, info, initMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	info = std.MockInfo("random", nil)
	handleMsg := []byte(`{"increment":{}}`)
	_, err = Handle(deps, env, info, handleMsg)
	require.NoError(t, err)

	qmsg := []byte(`{"get_count":{}}`)
	data, err := Query(deps, env, qmsg)
	require.NoError(t, err)
	var qres CountResponse
	err = json.Unmarshal(data.Ok, &qres)
	require.NoError(t, err)
	require.Equal(t, uint64(124), qres.Count)
}
