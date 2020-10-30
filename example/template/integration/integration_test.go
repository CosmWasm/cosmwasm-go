package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"

	mocks "github.com/CosmWasm/go-cosmwasm/api"
	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/example/TEMPLATE/src"
	"github.com/cosmwasm/cosmwasm-go/std/integration"
)

var CONTRACT = filepath.Join("..", "TEMPLATE.wasm")

func TestWorkflow(t *testing.T) {
	instance := integration.NewInstance(t, CONTRACT, 100_000_000, nil)

	env := mocks.MockEnv()
	info := mocks.MockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", nil)
	initMsg := []byte(`{"count":1234}`)
	res, _, err := instance.Init(env, info, initMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	// both work, let's test that
	//queryMsg := []byte(`{"get_count":{"b": "c"}}`)
	queryMsg := []byte(`{"get_count":{}}`)
	qres, _, err := instance.Query(env, queryMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	//require.Equal(t, uint64(0xb00c7), _)

	// let us parse the query??
	var count src.CountResponse
	err = json.Unmarshal(qres, &count)
	require.NoError(t, err)
	require.Equal(t, uint64(1234), count.Count)
}
