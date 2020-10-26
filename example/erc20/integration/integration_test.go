package integration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	cosmwasm "github.com/CosmWasm/go-cosmwasm"
	types "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/example/erc20/src"
)

var CONTRACT = filepath.Join("..", "erc20.wasm")

const FEATURES = "staking"
const mockContractAddr = "contract"

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

	// a whole lot of setup object using go-cosmwasm mock/test code
	var gasLimit uint64 = 100_000_000
	gasMeter := NewMockGasMeter(gasLimit)
	store := NewLookup(gasMeter)
	api := NewMockAPI()
	querier := DefaultQuerier(mockContractAddr, types.Coins{types.NewCoin(100, "ATOM")})
	env := mockEnv("creator")

	initMsg := []byte(`{"name":"OKB","symbol":"OKB","decimal":10,"total_supply":170000}`)
	_, gas, err := wasmer.Instantiate(codeID,
		env,
		initMsg,
		store,
		api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.Equal(t, uint64(486334), gas)
	require.NoError(t, err)

	handleMsg := []byte(`{"Transfer":{"to":"1234567","value": 2000}}`)
	_, gas, err = wasmer.Execute(codeID,
		env,
		handleMsg,
		store,
		api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.NoError(t, err)
	require.Equal(t, uint64(400000), gas)

	queryMsg := []byte(`{"balance":{"address":"1234567"}}`)
	qres, gas, err := wasmer.Query(codeID,
		queryMsg,
		store,
		api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.NoError(t, err)
	require.NotEmpty(t, qres)
	require.Equal(t, uint64(200000), gas)

	// let us parse the query??
	var bal src.BalanceResponse
	err = json.Unmarshal(qres, &bal)
	require.NoError(t, err)
	require.Equal(t, uint64(2000), bal.Value)

}

func mockEnv(sender types.HumanAddress) types.Env {
	return types.Env{
		Block: types.BlockInfo{
			Height:  123,
			Time:    1578939743,
			ChainID: "foobar",
		},
		Message: types.MessageInfo{
			Sender: sender,
			SentFunds: []types.Coin{{
				Denom:  "ATOM",
				Amount: "100",
			}},
		},
		Contract: types.ContractInfo{
			Address: mockContractAddr,
		},
	}
}
