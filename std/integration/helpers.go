package integration

import (
	"io/ioutil"
	"os"
	"testing"

	cosmwasm "github.com/CosmWasm/go-cosmwasm"
	"github.com/stretchr/testify/require"
)

const FEATURES = "staking"

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
