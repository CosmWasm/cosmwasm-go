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

	"github.com/cosmwasm/cosmwasm-go/example/hackatom/src"
)

var CONTRACT = filepath.Join("..", "hackatom.wasm")

const FEATURES = "staking"
const mockContractAddr = "coral1lstq3dy9v0s86czkx0rvgwnmunds5y2lz53all"

func setup(t *testing.T) (*cosmwasm.Wasmer, []byte) {
	// setup wasmer instance
	tmpdir, err := ioutil.TempDir("", "erc20")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpdir) })

	wasmer, err := cosmwasm.NewWasmer(tmpdir, FEATURES, true)
	require.NoError(t, err)

	// upload code and get some sha256 hash
	bz, err := ioutil.ReadFile(CONTRACT)
	require.NoError(t, err)
	codeID, err := wasmer.Create(bz)
	require.NoError(t, err)
	require.Equal(t, 32, len(codeID))

	return wasmer, codeID
}

func TestWorkflow(t *testing.T) {
	wasmer, codeID := setup(t)

	// a whole lot of setup object using go-cosmwasm mock/test code
	var gasLimit uint64 = 100_000_000
	gasMeter := NewMockGasMeter(gasLimit)
	store := NewLookup(gasMeter)
	api := NewMockAPI()
	querier := DefaultQuerier(mockContractAddr, nil)

	info := mockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", nil)
	initMsg := []byte(`{"count":1234}`)
	res, _, err := wasmer.Instantiate(codeID,
		mockEnv(),
		info,
		initMsg,
		store,
		api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.NoError(t, err)
	require.NotNil(t, res)

	queryMsg := []byte(`{"get_count":{"a":"b"}}`)
	qres, _, err := wasmer.Query(codeID,
		mockEnv(),
		queryMsg,
		store,
		api,
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

func mockEnv() types.Env {
	return types.Env{
		Block: types.BlockInfo{
			Height:  123,
			Time:    1578939743,
			ChainID: "foobar",
		},
		Contract: types.ContractInfo{
			Address: mockContractAddr,
		},
	}
}

func mockInfo(sender types.HumanAddress, sentFunds []types.Coin) types.MessageInfo {
	return types.MessageInfo{
		Sender:    sender,
		SentFunds: sentFunds,
	}
}
