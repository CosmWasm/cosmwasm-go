package src

import (
	"encoding/json"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func encode(t *testing.T, o json.Marshaler) []byte {
	b, err := o.MarshalJSON()
	require.NoError(t, err)
	return b
}

func TestContract(t *testing.T) {
	const exists = "existing"
	deps := mocks.Deps(nil)
	env := mocks.Env()
	info := mocks.Info(exists, nil)

	// create
	_, err := Execute(deps, env, info, encode(t, ExecuteMsg{CreateIdentity: &MsgCreateIdentity{
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Shiganshina District",
		PostalCode: 100,
	}}))
	require.NoError(t, err)

	// assert it exists
	identityReq := encode(t, QueryMsg{Identity: &QueryIdentity{ID: exists}})

	resp, err := Query(deps, env, identityReq)
	require.NoError(t, err)

	identity := new(Person)
	require.NoError(t, identity.UnmarshalJSON(resp))

	require.Equal(t, Person{
		Address:    exists,
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Shiganshina District",
		PostalCode: 100,
	}, *identity)

	// already exists
	_, err = Execute(deps, env, info, encode(t, ExecuteMsg{CreateIdentity: &MsgCreateIdentity{
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Shiganshina District",
		PostalCode: 100,
	}}))
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPersonAlreadyExists))

	// update
	_, err = Execute(deps, env, info, encode(t, &ExecuteMsg{UpdateCity: &MsgUpdateCity{
		City:       "Liberio",
		PostalCode: 200,
	}}))
	require.NoError(t, err)

	// check if update went fine
	resp, err = Query(deps, env, identityReq)
	require.NoError(t, err)

	identity = new(Person)
	require.NoError(t, identity.UnmarshalJSON(resp))

	require.Equal(t, Person{
		Address:    exists,
		Name:       "Eren",
		Surname:    "Yeager",
		City:       "Liberio",
		PostalCode: 200,
	}, *identity)

	// delete
	_, err = Execute(deps, env, info, encode(t, ExecuteMsg{DeleteIdentity: &MsgDelete{}}))
	require.NoError(t, err)

	// we test correct not found handling

	// read
	_, err = Query(deps, env, identityReq)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPersonNotFound))
	// update
	_, err = Execute(deps, env, info, encode(t, ExecuteMsg{UpdateCity: &MsgUpdateCity{City: "milan", PostalCode: 1}}))
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPersonNotFound))
	// delete
	_, err = Execute(deps, env, info, encode(t, ExecuteMsg{DeleteIdentity: &MsgDelete{}})) // delete
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPersonNotFound))
}
