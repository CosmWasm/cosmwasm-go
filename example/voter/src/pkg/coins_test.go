package pkg

import (
	"testing"

	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
	"github.com/stretchr/testify/assert"
)

func TestCoinsContainMinAmount(t *testing.T) {
	type testCase struct {
		name  string
		coins []stdTypes.Coin
		coin  stdTypes.Coin
		//
		errExpected bool
	}

	testCases := []testCase{
		{
			name: "OK: GT",
			coins: []stdTypes.Coin{
				stdTypes.NewCoinFromUint64(100, "uatom"),
				stdTypes.NewCoinFromUint64(50, "musdt"),
			},
			coin: stdTypes.NewCoinFromUint64(45, "musdt"),
		},
		{
			name: "OK: GT",
			coins: []stdTypes.Coin{
				stdTypes.NewCoinFromUint64(100, "uatom"),
				stdTypes.NewCoinFromUint64(50, "musdt"),
			},
			coin: stdTypes.NewCoinFromUint64(50, "musdt"),
		},
		{
			name: "Fail: not found",
			coins: []stdTypes.Coin{
				stdTypes.NewCoinFromUint64(100, "uatom"),
				stdTypes.NewCoinFromUint64(50, "musdt"),
			},
			coin:        stdTypes.NewCoinFromUint64(10, "uusdc"),
			errExpected: true,
		},
		{
			name: "Fail: LT",
			coins: []stdTypes.Coin{
				stdTypes.NewCoinFromUint64(100, "uatom"),
				stdTypes.NewCoinFromUint64(50, "musdt"),
			},
			coin:        stdTypes.NewCoinFromUint64(55, "musdt"),
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errExpected {
				assert.Error(t, CoinsContainMinAmount(tc.coins, tc.coin))
				return
			}
			assert.NoError(t, CoinsContainMinAmount(tc.coins, tc.coin))
		})
	}
}
