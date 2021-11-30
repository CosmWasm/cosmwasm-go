package integration

import (
	"encoding/json"
	mocks "github.com/CosmWasm/wasmvm/api"
	"github.com/cosmwasm/cosmwasm-go/example/queue/src"
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
	// queue empty
	_, gas, err := i.Execute(env, info, encode(t, src.ExecuteMsg{Enqueue: &src.Enqueue{Value: 5}}))
	require.NoError(t, err)
	t.Logf("empty queue cost: %d", gas)
	// check one item in the queue only
	respBytes, _, err := i.Query(env, encode(t, src.QueryMsg{Count: &struct{}{}}))
	require.NoError(t, err)
	resp := new(src.CountResponse)
	require.NoError(t, resp.UnmarshalJSON(respBytes))
	require.Equal(t, resp.Count, uint32(1))

	// enqueue another item
	_, _, err = i.Execute(env, info, encode(t, src.ExecuteMsg{Enqueue: &src.Enqueue{Value: 6}}))
	require.NoError(t, err)
	// dequeue two items so we can check if results are the expected ones
	exResp, gas, err := i.Execute(env, info, encode(t, src.ExecuteMsg{Dequeue: &src.Dequeue{}}))
	require.NoError(t, err)
	require.NotEmpty(t, exResp.Data)
	item := new(src.Item)
	require.NoError(t, item.UnmarshalJSON(exResp.Data))
	require.Equal(t, int32(5), item.Value)

	exResp, gas, err = i.Execute(env, info, encode(t, src.ExecuteMsg{Dequeue: &src.Dequeue{}}))
	require.NoError(t, err)
	require.NotEmpty(t, exResp.Data)
	item = new(src.Item)
	require.NoError(t, item.UnmarshalJSON(exResp.Data))
	require.Equal(t, int32(6), item.Value)
}
