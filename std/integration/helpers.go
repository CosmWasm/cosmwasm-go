package integration

import (
	"io/ioutil"
	"os"
	"testing"

	cosmwasm "github.com/CosmWasm/go-cosmwasm"
	mocks "github.com/CosmWasm/go-cosmwasm/api"
	types "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/stretchr/testify/require"
)

const FEATURES = "staking"

// TODO: move this into go-cosmwasm at some point
func NewCoins(amount uint64, denom string) []types.Coin {
	return []types.Coin{types.NewCoin(amount, denom)}
}

// End transient code

func SetupWasmer(t *testing.T, contractPath string) (*cosmwasm.Wasmer, []byte) {
	// setup wasmer instance
	tmpdir, err := ioutil.TempDir("", "wasmer")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpdir) })

	wasmer, err := cosmwasm.NewWasmer(tmpdir, FEATURES, true)
	require.NoError(t, err)
	codeID := StoreCode(t, wasmer, contractPath)

	return wasmer, codeID
}

// Returns code id
func StoreCode(t *testing.T, wasmer *cosmwasm.Wasmer, contractPath string) []byte {
	// upload code and get some sha256 hash
	bz, err := ioutil.ReadFile(contractPath)
	require.NoError(t, err)
	codeID, err := wasmer.Create(bz)
	require.NoError(t, err)
	require.Equal(t, 32, len(codeID))
	return codeID
}

type Instance struct {
	Wasmer   *cosmwasm.Wasmer
	CodeID   []byte
	GasLimit uint64
	GasMeter cosmwasm.GasMeter
	Store    *mocks.Lookup
	Api      *mocks.GoAPI
	Querier  mocks.Querier
}

func NewInstance(t *testing.T, contractPath string, gasLimit uint64, funds []types.Coin) Instance {
	wasmer, codeID := SetupWasmer(t, contractPath)
	gasMeter := mocks.NewMockGasMeter(gasLimit)

	return Instance{
		Wasmer:   wasmer,
		CodeID:   codeID,
		GasLimit: gasLimit,
		GasMeter: gasMeter,
		Store:    mocks.NewLookup(gasMeter),
		Api:      mocks.NewMockAPI(),
		Querier:  mocks.DefaultQuerier(mocks.MOCK_CONTRACT_ADDR, funds),
	}
}

func (i *Instance) Init(env types.Env, info types.MessageInfo, initMsg []byte) (*types.InitResponse, uint64, error) {
	return i.Wasmer.Instantiate(
		i.CodeID,
		env,
		info,
		initMsg,
		i.Store,
		*i.Api,
		i.Querier,
		i.GasMeter,
		i.GasLimit,
	)
}

func (i *Instance) Handle(env types.Env, info types.MessageInfo, handleMsg []byte) (*types.HandleResponse, uint64, error) {
	return i.Wasmer.Execute(
		i.CodeID,
		env,
		info,
		handleMsg,
		i.Store,
		*i.Api,
		i.Querier,
		i.GasMeter,
		i.GasLimit,
	)
}

func (i *Instance) Query(env types.Env, queryMsg []byte) ([]byte, uint64, error) {
	return i.Wasmer.Query(
		i.CodeID,
		env,
		queryMsg,
		i.Store,
		*i.Api,
		i.Querier,
		i.GasMeter,
		i.GasLimit,
	)
}
