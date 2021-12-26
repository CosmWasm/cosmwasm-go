package src

import (
	"github.com/cosmwasm/cosmwasm-go/std/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContract_ExecCreateIdentity(t *testing.T) {
	deps := mock.Deps(nil)
	env := mock.Env()
	info := mock.Info("sender", nil)

	// test create
	_, err := Contract{}.ExecCreateIdentity(deps, &env, &info, &MsgCreateIdentity{
		Name:       "name",
		Surname:    "surname",
		City:       "city",
		PostalCode: 100,
	})
	require.NoError(t, err)

	// check if it exists
	person, err := Contract{}.QueryIdentity(deps, &env, &QueryIdentity{
		ID: info.Sender,
	})
	require.NoError(t, err)
	require.Equal(t, Person{
		Address:    info.Sender,
		Name:       "name",
		Surname:    "surname",
		City:       "city",
		PostalCode: 100,
	}, *person)

	// check can create only once
	_, err = Contract{}.ExecCreateIdentity(deps, &env, &info, &MsgCreateIdentity{
		Name:       "name",
		Surname:    "surname",
		City:       "city",
		PostalCode: 100,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrPersonAlreadyExists.Error())
}

func TestContract_ExecDeleteIdentity(t *testing.T) {
	deps := mock.Deps(nil)
	env := mock.Env()
	info := mock.Info("sender", nil)

	// create
	_, err := Contract{}.ExecCreateIdentity(deps, &env, &info, &MsgCreateIdentity{
		Name:       "name",
		Surname:    "surname",
		City:       "city",
		PostalCode: 100,
	})
	require.NoError(t, err)
	// delete
	_, err = Contract{}.ExecDeleteIdentity(deps, &env, &info, &MsgDelete{})
	require.NoError(t, err)

	// test can't delete
	_, err = Contract{}.QueryIdentity(deps, &env, &QueryIdentity{ID: info.Sender})
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrPersonNotFound.Error())
}

func TestContract_ExecUpdateCity(t *testing.T) {
	deps := mock.Deps(nil)
	env := mock.Env()
	info := mock.Info("sender", nil)

	// test error if person does not exist
	_, err := Contract{}.ExecUpdateCity(deps, &env, &info, &MsgUpdateCity{
		City:       "othercity",
		PostalCode: 300,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrPersonNotFound.Error())

	// create
	_, err = Contract{}.ExecCreateIdentity(deps, &env, &info, &MsgCreateIdentity{
		Name:       "name",
		Surname:    "surname",
		City:       "city",
		PostalCode: 100,
	})
	require.NoError(t, err)

	// update
	_, err = Contract{}.ExecUpdateCity(deps, &env, &info, &MsgUpdateCity{
		City:       "othercity",
		PostalCode: 1000,
	})
	require.NoError(t, err)

	person, err := Contract{}.QueryIdentity(deps, &env, &QueryIdentity{ID: info.Sender})
	require.NoError(t, err)

	require.Equal(t, *person, Person{
		Address:    info.Sender,
		Name:       "name",
		Surname:    "surname",
		City:       "othercity",
		PostalCode: 1000,
	})
}
