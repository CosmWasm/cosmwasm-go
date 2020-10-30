package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"

	mocks "github.com/CosmWasm/go-cosmwasm/api"
	"github.com/CosmWasm/go-cosmwasm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/example/hackatom/src"
	"github.com/cosmwasm/cosmwasm-go/std/integration"
)

var CONTRACT = filepath.Join("..", "hackatom.wasm")

func mustEncode(t *testing.T, msg interface{}) []byte {
	bz, err := json.Marshal(msg)
	require.NoError(t, err)
	return bz
}

const VERIFIER = "verifies"
const BENEFICIARY = "benefits"
const FUNDER = "creator"

// this can be used for a quick setup if you don't have nay other requirements
func defaultInit(t *testing.T, funds []types.Coin) *integration.Instance {
	instance := integration.NewInstance(t, CONTRACT, 100_000_000, funds)

	env := mocks.MockEnv()
	info := mocks.MockInfo(FUNDER, funds)
	initMsg := src.InitMsg{
		Verifier:    VERIFIER,
		Beneficiary: BENEFICIARY,
	}
	res, _, err := instance.Init(env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	return &instance
}

func TestInitAndQuery(t *testing.T) {
	instance := integration.NewInstance(t, CONTRACT, 100_000_000, nil)

	env := mocks.MockEnv()
	info := mocks.MockInfo(FUNDER, nil)
	initMsg := src.InitMsg{
		Verifier:    VERIFIER,
		Beneficiary: BENEFICIARY,
	}
	res, _, err := instance.Init(env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 0, len(res.Messages))
	assert.Equal(t, 1, len(res.Attributes))
	attr := res.Attributes[0]
	assert.Equal(t, "Let the", attr.Key)
	assert.Equal(t, "hacking begin", attr.Value)

	qmsg := []byte(`{"verifier":{}}`)
	data, _, err := instance.Query(env, qmsg)
	require.NoError(t, err)
	var qres src.VerifierResponse
	err = json.Unmarshal(data, &qres)
	require.NoError(t, err)
	assert.Equal(t, VERIFIER, qres.Verifier)
}

func TestPanic(t *testing.T) {
	deps := defaultInit(t, nil)
	env := mocks.MockEnv()
	info := mocks.MockInfo(FUNDER, nil)
	handleMsg := []byte(`{"panic":{}}`)
	_, _, err := deps.Handle(env, info, handleMsg)
	require.Error(t, err)
}

//
//func TestWorkflow(t *testing.T) {
//	instance := integration.NewInstance(t, CONTRACT, 100_000_000, nil)
//
//	env := mocks.MockEnv()
//	info := mocks.MockInfo("coral1e86v774dch5uwkks0cepw8mdz8a9flhhapvf6w", nil)
//	initMsg := []byte(`{"count":1234}`)
//	res, _, err := instance.Init(env, info, initMsg)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	// both work, let's test that
//	//queryMsg := []byte(`{"get_count":{"b": "c"}}`)
//	queryMsg := []byte(`{"get_count":{}}`)
//	qres, _, err := instance.Query(env, queryMsg)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//	//require.Equal(t, uint64(0xb00c7), _)
//
//	// let us parse the query??
//	var count src.CountResponse
//	err = json.Unmarshal(qres, &count)
//	require.NoError(t, err)
//	require.Equal(t, uint64(1234), count.Count)
//}
