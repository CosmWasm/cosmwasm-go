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

func TestRelease(t *testing.T) {
	cases := map[string]struct {
		signer string
		funds  []types.Coin
		valid  bool
	}{
		"verifier releases": {VERIFIER, integration.NewCoins(765432, "wei"), true},
		"random fails":      {BENEFICIARY, integration.NewCoins(123456, "wei"), false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			deps := defaultInit(t, tc.funds)
			env := mocks.MockEnv()
			info := mocks.MockInfo(tc.signer, nil)
			handleMsg := []byte(`{"release":{}}`)
			res, _, err := deps.Handle(env, info, handleMsg)
			if !tc.valid {
				require.Error(t, err)
				require.Equal(t, "Handle: Unauthorized", err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Equal(t, 1, len(res.Messages))
				msg := res.Messages[0]
				expected := types.CosmosMsg{Bank: &types.BankMsg{Send: &types.SendMsg{
					FromAddress: mocks.MOCK_CONTRACT_ADDR,
					ToAddress:   BENEFICIARY,
					Amount:      tc.funds,
				}}}
				assert.Equal(t, expected, msg)
				assert.Equal(t, 2, len(res.Attributes))
				assert.Equal(t, []types.EventAttribute{{"action", "release"}, {"destination", BENEFICIARY}}, res.Attributes)
			}
		})
	}
}

func TestQueryOther(t *testing.T) {
	contractFunds := []types.Coin{
		types.NewCoin(1000, "wei"),
		types.NewCoin(555, "uatom"),
	}
	deps := defaultInit(t, contractFunds)
	env := mocks.MockEnv()

	// TODO: set some balances
	richFunds := []types.Coin{
		types.NewCoin(123456789, "uatom"),
		types.NewCoin(9876542, "satoshi"),
		types.NewCoin(557755, "utgd"),
	}
	deps.SetQuerierBalance("rich", richFunds)

	cases := map[string]struct {
		account string
		balance []types.Coin
	}{
		"contract self": {mocks.MOCK_CONTRACT_ADDR, contractFunds},
		"random":        {"random", nil},
		"rich":          {"rich", richFunds},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// json encoding makes invalid QueryMsg... look into this later (only ezjson works?)
			queryMsg := []byte(`{"other_balance":{"address":"` + tc.account + `"}}`)

			raw, _, err := deps.Query(env, queryMsg)
			require.NoError(t, err)
			var res types.AllBalancesResponse
			err = json.Unmarshal(raw, &res)
			require.NoError(t, err)
			require.Equal(t, types.Coins(tc.balance), res.Amount)
		})
	}
}

//let rich_addr = HumanAddr::from("foobar");
//let rich_balance = coins(10000, "gold");
//let deps = mock_dependencies_with_balances(&[(&rich_addr, &rich_balance)]);
