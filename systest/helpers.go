package systest

import (
	"io/ioutil"
	"os"
	"testing"

	wasmvm "github.com/CosmWasm/wasmvm"
	mocks "github.com/CosmWasm/wasmvm/api"
	types "github.com/CosmWasm/wasmvm/types"

	unitmocks "github.com/cosmwasm/cosmwasm-go/std/mock"

	"github.com/stretchr/testify/require"
)

const FEATURES = "staking"

// TODO: move this into wasmvm at some point
func NewCoins(amount uint64, denom string) []types.Coin {
	return []types.Coin{types.NewCoin(amount, denom)}
}

var deserCost = types.UFraction{
	Numerator:   1,
	Denominator: 10,
}

// End transient code

func SetupWasmer(t *testing.T, contractPath string) (*wasmvm.VM, []byte) {
	// setup wasmer instance
	tmpdir, err := ioutil.TempDir("", "wasmer")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpdir) })

	wasmer, err := wasmvm.NewVM(tmpdir, FEATURES, 256, true, 0)
	require.NoError(t, err)
	codeID := StoreCode(t, wasmer, contractPath)

	return wasmer, codeID
}

// Returns code id
func StoreCode(t *testing.T, wasmer *wasmvm.VM, contractPath string) []byte {
	// upload code and get some sha256 hash
	bz, err := ioutil.ReadFile(contractPath)
	require.NoError(t, err)
	codeID, err := wasmer.Create(bz)
	require.NoError(t, err)
	require.Equal(t, 32, len(codeID))
	return codeID
}

type Instance struct {
	Wasmer   *wasmvm.VM
	CodeID   []byte
	GasLimit uint64
	GasMeter wasmvm.GasMeter
	Store    *mocks.Lookup
	Api      *mocks.GoAPI
	Querier  mocks.Querier
}

func (i *Instance) SetQuerierBalance(addr string, balance []types.Coin) {
	mq := i.Querier.(mocks.MockQuerier)
	mq.Bank.Balances[addr] = balance
}

func NewInstance(t *testing.T, contractPath string, gasLimit uint64, funds []types.Coin) Instance {
	wasmer, codeID := SetupWasmer(t, contractPath)
	gasMeter := mocks.NewMockGasMeter(gasLimit)

	// we use callbacks that use the same logic as our unit tests
	mockApi := mocks.NewMockAPI()
	mockApi.HumanAddress = func(canon []byte) (string, uint64, error) {
		human, err := unitmocks.API{}.HumanAddress(canon)
		return human, 5000, err
	}
	mockApi.CanonicalAddress = func(human string) ([]byte, uint64, error) {
		canon, err := unitmocks.API{}.CanonicalAddress(human)
		return canon, 4000, err
	}

	return Instance{
		Wasmer:   wasmer,
		CodeID:   codeID,
		GasLimit: gasLimit,
		GasMeter: gasMeter,
		Store:    mocks.NewLookup(gasMeter),
		Api:      mockApi,
		Querier:  mocks.DefaultQuerier(mocks.MOCK_CONTRACT_ADDR, funds),
	}
}

func (i *Instance) Instantiate(env types.Env, info types.MessageInfo, initMsg []byte) (*types.Response, uint64, error) {
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
		deserCost,
	)
}

func (i *Instance) Execute(env types.Env, info types.MessageInfo, handleMsg []byte) (*types.Response, uint64, error) {
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
		deserCost,
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
		deserCost,
	)
}

func (i *Instance) Migrate(env types.Env, migrateMsg []byte) (*types.Response, uint64, error) {
	return i.Wasmer.Migrate(
		i.CodeID,
		env,
		migrateMsg,
		i.Store,
		*i.Api,
		i.Querier,
		i.GasMeter,
		i.GasLimit,
		deserCost,
	)
}
