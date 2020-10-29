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
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

var CONTRACT = filepath.Join("..", "erc20.wasm")

const FEATURES = "staking"
const mockContractAddr = "coral1lstq3dy9v0s86czkx0rvgwnmunds5y2lz53all"

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
	wasmer, err := cosmwasm.NewWasmer(tmpdir, FEATURES, true)
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
	info := mockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", nil)

	initMsg := []byte(`{"name":"OKB","symbol":"OKB","decimal":10,"total_supply":170000}`)
	//initMsg := []byte(`{123]]`) // invalid json
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
	//require.Equal(t, uint64(0xb787e), _)

	// check we get the attributes out
	require.Equal(t, 1, len(res.Attributes))
	require.Equal(t, "hello", res.Attributes[0].Key)
	require.Equal(t, "world", res.Attributes[0].Value)

	handleMsg := []byte(`{"Transfer":{"to":"1234567","value": 2000}}`)
	_, _, err = wasmer.Execute(codeID,
		mockEnv(),
		info,
		handleMsg,
		store,
		api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.NoError(t, err)
	//require.Equal(t, uint64(0x195920), _)

	queryMsg := []byte(`{"balance":{"address":"1234567"}}`)
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
	var bal src.BalanceResponse
	err = json.Unmarshal(qres, &bal)
	require.NoError(t, err)
	require.Equal(t, uint64(2000), bal.Value)
}

func TestInfoMarshalCompatibility(t *testing.T) {
	cases := map[string]struct {
		funds []types.Coin
	}{
		"no funds":   {},
		"some funds": {funds: types.Coins{types.NewCoin(1000, "uatom"), types.NewCoin(5000, "eth")}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			info := mockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", tc.funds)
			bz, err := json.Marshal(info)
			require.NoError(t, err)

			var parsed std.MessageInfo
			err = ezjson.Unmarshal(bz, &parsed)
			require.NoError(t, err)

			// types are different to compare, but re-encode should match
			reencode, err := ezjson.Marshal(parsed)
			require.NoError(t, err)
			require.Equal(t, string(bz), string(reencode))
		})
	}
}

func TestErrorReturned(t *testing.T) {
	// setup wasmer instance
	tmpdir, err := ioutil.TempDir("", "erc20")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	wasmer, err := cosmwasm.NewWasmer(tmpdir, FEATURES, true)
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
	querier := DefaultQuerier(mockContractAddr, nil)
	info := mockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", nil)

	// this doesn't create json parse error, but we error on validating an empty InitMsg
	initMsg := []byte(`{.32r!`)
	_, _, err = wasmer.Instantiate(codeID,
		mockEnv(),
		info,
		initMsg,
		store,
		api,
		querier,
		gasMeter,
		gasLimit,
	)
	require.Error(t, err)
	require.Equal(t, "Init: Name must be at least 2 characters", err.Error())

}

func TestWorkflowWithFunds(t *testing.T) {
	// setup wasmer instance
	tmpdir, err := ioutil.TempDir("", "erc20")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	wasmer, err := cosmwasm.NewWasmer(tmpdir, FEATURES, true)
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
	funds := types.Coins{types.NewCoin(1000, "uatom"), types.NewCoin(60000, "utgd")}
	querier := DefaultQuerier(mockContractAddr, funds)
	info := mockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", funds)

	initMsg := []byte(`{"name":"Okay","symbol":"OKB","decimal":10,"total_supply":170000}`)
	res, gas, err := wasmer.Instantiate(codeID,
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
	require.Equal(t, uint64(0xe4c9e), gas)
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
