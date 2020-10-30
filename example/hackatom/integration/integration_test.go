package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"

	cosmwasm "github.com/CosmWasm/go-cosmwasm/api"
	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/example/hackatom/src"
	"github.com/cosmwasm/cosmwasm-go/std/integration"
)

var CONTRACT = filepath.Join("..", "hackatom.wasm")

func TestWorkflow(t *testing.T) {
	wasmer, codeID := integration.SetupWasmer(t, CONTRACT)

	// a whole lot of setup object using go-cosmwasm mock/test code
	var gasLimit uint64 = 100_000_000
	gasMeter := cosmwasm.NewMockGasMeter(gasLimit)
	store := cosmwasm.NewLookup(gasMeter)
	api := cosmwasm.NewMockAPI()
	querier := cosmwasm.DefaultQuerier(cosmwasm.MOCK_CONTRACT_ADDR, nil)

	info := cosmwasm.MockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", nil)
	initMsg := []byte(`{"count":1234}`)
	res, _, err := wasmer.Instantiate(codeID,
		cosmwasm.MockEnv(),
		info,
		initMsg,
		store,
		*api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.NoError(t, err)
	require.NotNil(t, res)

	//queryMsg := []byte(`{"get_count":{"b": "c"}}`)
	queryMsg := []byte(`{"get_count":{}}`)
	qres, _, err := wasmer.Query(codeID,
		cosmwasm.MockEnv(),
		queryMsg,
		store,
		*api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.NoError(t, err)
	require.NotEmpty(t, qres)
	//require.Equal(t, uint64(0xb00c7), _)

	// let us parse the query??
	var count src.CountResponse
	err = json.Unmarshal(qres, &count)
	require.NoError(t, err)
	require.Equal(t, uint64(1234), count.Count)
}
