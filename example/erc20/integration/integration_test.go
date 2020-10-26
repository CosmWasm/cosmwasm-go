package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	cosmwasm "github.com/CosmWasm/go-cosmwasm"
)

var CONTRACT = filepath.Join("..", "erc20.wasm")
const FEATURES = "staking"

func loadCode(t *testing.T) []byte {
	bz, err := ioutil.ReadFile(CONTRACT)
	require.NoError(t, err)
	return bz
}

func TestWorkflow(t *testing.T) {
	// setup wasmer instance
	tmpdir, err := ioutil.TempDir("", "erc20")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	wasmer, err := cosmwasm.NewWasmer(tmpdir, FEATURES, 0)
	require.NoError(t, err)

	// upload code and get some sha256 hash
	codeID, err := wasmer.Create(loadCode(t))
	require.NoError(t, err)
	require.Equal(t, 32, len(codeID))

	//deps := std.MockExtern()
	//env := std.MockEnv("original_owner_addr", nil)
	//
	//// TODO: try to init it
	//initMsg := []byte(`{"name":"OKB","symbol":"OKB","decimal":10,"total_supply":170000}`)
	//res, gas, err := wasmer.Instantiate(codeID,
	//	env types.Env,
	//	initMsg,
	//	store KVStore,
	//	goapi GoAPI,
	//	querier Querier,
	//	gasMeter GasMeter,
	//	gasLimit uint64,
	//)
	//require.NoError(t, err)
	//
	//handleMsg := []byte(`{"Transfer":{"to":"1234567","value": 2000}}`)
	//res, gas, err := wasmer.Execute(codeID,
	//	env types.Env,
	//	handleMsg,
	//	store KVStore,
	//	goapi GoAPI,
	//	querier Querier,
	//	gasMeter GasMeter,
	//	gasLimit uint64,
	//)
	//require.NoError(t, err)
	//
	//queryMsg := []byte(`{"balance":{"address":"1234567"}}`)
	//qres, gas, err := wasmer.Query(codeID,
	//	queryMsg,
	//	store KVStore,
	//	goapi GoAPI,
	//	querier Querier,
	//	gasMeter GasMeter,
	//	gasLimit uint64,
	//)
	//require.NoError(t, err)
	//require.NotEmpty(t, qres)
	//
	//// let us parse the query??
	//var bal src.BalanceResponse
	//err = json.Unmarshal(qres, &bal)
	//require.NoError(t, err)
	//require.Equal(t, uint64(2000), bal.Value)

}
