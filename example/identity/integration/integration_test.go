package integration

import (
	"encoding/json"
	mocks "github.com/CosmWasm/wasmvm/api"
	"github.com/cosmwasm/cosmwasm-go/example/identity/src"
	"github.com/cosmwasm/cosmwasm-go/systest"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

var contractPath = filepath.Join("..", "queue.wasm")

func encode(t *testing.T, o json.Marshaler) []byte {
	bytes, err := o.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

func instance(t *testing.T) *systest.Instance {
	i := systest.NewInstance(t, contractPath, 15_000_000_000_000, nil)
	return &i
}

func TestExecute(t *testing.T) {
	i := instance(t)
	env := mocks.MockEnv()
	info := mocks.MockInfo("none", nil)

	// test create
	_, gas, err := i.Execute(env, info, encode(t, &src.ExecuteMsg{CreateIdentity: &src.MsgCreateIdentity{
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Shiganshina District",
		PostalCode: 100,
	}}))

	require.NoError(t, err)
	t.Logf("create gas: %d", gas)

	// test read
	// resp, gas, err := i.Query(env, encode())
}
