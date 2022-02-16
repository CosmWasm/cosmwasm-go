package src

import (
	"github.com/cosmwasm/cosmwasm-go/std/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var contract = Contract{}

func TestContract(t *testing.T) {
	deps := mock.Deps(nil)
	info := mock.Info("sender", nil)
	env := mock.Env()

	// test creation
	_, err := contract.CreateIdentity(deps, &env, &info, &MsgCreateIdentity{
		Name:       "Name",
		Surname:    "Surname",
		City:       "City",
		PostalCode: 100,
	})
	require.NoError(t, err)

	// test query
	person, err := contract.QueryIdentity(deps, &env, &QueryIdentity{ID: info.Sender})
	require.NoError(t, err)

	require.Equal(t, person.Person, &Person{
		Address:    info.Sender,
		Name:       "Name",
		Surname:    "Surname",
		City:       "City",
		PostalCode: 100,
	})

	// test update
	_, err = contract.UpdateCity(deps, &env, &info, &MsgUpdateCity{
		City:       "UpdatedCity",
		PostalCode: 200,
	})
	require.NoError(t, err)
	person, err = contract.QueryIdentity(deps, &env, &QueryIdentity{ID: info.Sender})
	require.NoError(t, err)
	require.Equal(t, person.Person, &Person{
		Address:    info.Sender,
		Name:       "Name",
		Surname:    "Surname",
		City:       "UpdatedCity",
		PostalCode: 200,
	})

	// test delete
	_, err = contract.DeleteIdentity(deps, &env, &info, &MsgDelete{})
	require.NoError(t, err)
	_, err = contract.QueryIdentity(deps, &env, &QueryIdentity{ID: info.Sender})
	require.Contains(t, err.Error(), ErrPersonNotFound.Error())
}
