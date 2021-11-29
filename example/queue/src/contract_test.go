package src

import (
	"encoding/json"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/mocks"
	"github.com/cosmwasm/cosmwasm-go/std/types"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func encode(t *testing.T, msg json.Marshaler) []byte {
	x, err := msg.MarshalJSON()
	require.NoError(t, err)
	return x
}

func Test_executeEnqueue(t *testing.T) {
	deps := mocks.MockDeps(nil)
	env := mocks.MockEnv()
	info := mocks.MockInfo("none", nil)
	// queue empty
	_, err := Execute(deps, env, info, encode(t, ExecuteMsg{Enqueue: &Enqueue{Value: 5}}))
	require.NoError(t, err)
	iter := deps.Storage.Range(nil, nil, std.Ascending)
	k, v, err := iter.Next()
	require.NoError(t, err)
	require.Equal(t, k, []byte{FirstKey})
	got := new(Item)
	require.NoError(t, got.UnmarshalJSON(v))
	// matching values
	require.Equal(t, &Item{Value: 5}, got)
	// no more results
	_, _, err = iter.Next()
	require.Error(t, err)
	require.Equal(t, err.Error(), std.ErrIteratorDone.Error())

	// queue not empty
	_, err = Execute(deps, env, info, encode(t, ExecuteMsg{Enqueue: &Enqueue{Value: 8}}))
	require.NoError(t, err)
	iter = deps.Storage.Range(nil, nil, std.Ascending)
	_, _, err = iter.Next() // skip first key
	require.NoError(t, err)
	k, v, err = iter.Next() // second key
	require.NoError(t, err)
	require.Equal(t, FirstKey+1, k[0]) // first key +1
	got = new(Item)
	require.NoError(t, got.UnmarshalJSON(v))
	// matching values
	require.Equal(t, &Item{Value: 8}, got)
	// no more results
	_, _, err = iter.Next()
	require.Error(t, err)
	require.Equal(t, err.Error(), std.ErrIteratorDone.Error())
}

func TestMigrate(t *testing.T) {
	deps := mocks.MockDeps(nil)

	_, err := executeEnqueue(deps, types.Env{}, types.MessageInfo{}, &Enqueue{Value: 1})
	require.NoError(t, err)

	_, err = Migrate(deps, types.Env{}, nil)
	require.NoError(t, err)

	iter := deps.Storage.Range(nil, nil, std.Ascending)
	for i := 100; i < 103; i++ {
		k, v, err := iter.Next()
		require.NoError(t, err)
		require.Equal(t, k[0], uint8(i-100))

		item := new(Item)
		require.NoError(t, item.UnmarshalJSON(v))
		require.Equal(t, item.Value, int32(i))
	}
}

func TestQuery_Count(t *testing.T) {
	deps := mocks.MockDeps(nil)
	env := mocks.MockEnv()
	info := mocks.MockInfo("none", nil)

	for i := int32(0); i < 20; i++ {
		_, err := executeEnqueue(deps, env, info, &Enqueue{Value: i})
		require.NoError(t, err)
	}

	respBytes, err := Query(deps, env, encode(t, QueryMsg{Count: &struct{}{}}))
	require.NoError(t, err)

	resp := new(CountResponse)
	require.NoError(t, resp.UnmarshalJSON(respBytes))

	require.Equal(t, uint32(20), resp.Count)
}

func TestQuery_Sum(t *testing.T) {
	deps := mocks.MockDeps(nil)
	env := mocks.MockEnv()
	info := mocks.MockInfo("none", nil)
	rand.Seed(time.Now().UnixNano())

	total := int32(0)
	for i := 0; i < rand.Intn(100); i++ {
		v := rand.Int31()
		_, err := executeEnqueue(deps, env, info, &Enqueue{Value: v})
		require.NoError(t, err)
		total += v
	}

	respBytes, err := Query(deps, env, encode(t, QueryMsg{Sum: &struct{}{}}))
	require.NoError(t, err)

	resp := new(SumResponse)
	require.NoError(t, resp.UnmarshalJSON(respBytes))

	require.Equal(t, total, resp.Sum)
}

func TestQuery_List(t *testing.T) {

}
