package integration

import (
	mocks "github.com/CosmWasm/wasmvm/api"
	"github.com/cosmwasm/cosmwasm-go/example/identityv2/src"
	"github.com/cosmwasm/cosmwasm-go/systest"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

var contractPath = filepath.Join("..", "identityv2.wasm")

func instance(t *testing.T) *systest.Instance {
	i := systest.NewInstance(t, contractPath, 15_000_000_000_000, nil)
	return &i
}

func TestExecute(t *testing.T) {
	i := instance(t)
	env := mocks.MockEnv()
	info := mocks.MockInfo("none", nil)

	// test create
	_, gas, err := i.Execute(env, info, (&src.MsgCreateIdentity{
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Shiganshina District",
		PostalCode: 100,
	}).AsExecuteMsg())

	require.NoError(t, err)
	t.Logf("create gas: %d", gas)

	// test read
	resp, gas, err := i.Query(env, (&src.QueryIdentity{ID: info.Sender}).AsQueryMsg())
	require.NoError(t, err)

	identityResp := new(src.Person)
	require.NoError(t, identityResp.UnmarshalJSON(resp))

	require.Equal(t, src.Person{
		Address:    info.Sender,
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Shiganshina District",
		PostalCode: 100,
	}, *identityResp)

	// test update
	_, gas, err = i.Execute(env, info, (&src.MsgUpdateCity{
		City:       "Liberio",
		PostalCode: 200,
	}).AsExecuteMsg())
	require.NoError(t, err)
	t.Logf("update gas: %d", gas)

	// check if update went fine
	resp, gas, err = i.Query(env, (&src.QueryIdentity{ID: info.Sender}).AsQueryMsg())
	require.NoError(t, err)

	identityResp = new(src.Person)
	require.NoError(t, identityResp.UnmarshalJSON(resp))

	require.Equal(t, src.Person{
		Address:    info.Sender,
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Liberio",
		PostalCode: 200,
	}, *identityResp)

	// delete
	_, gas, err = i.Execute(env, info, (&src.MsgDelete{}).AsExecuteMsg())
	require.NoError(t, err)

	// check deletion
	_, gas, err = i.Query(env, (&src.QueryIdentity{ID: info.Sender}).AsQueryMsg())
	require.Error(t, err)
	require.Contains(t, err.Error(), src.ErrPersonNotFound.Error())
}
