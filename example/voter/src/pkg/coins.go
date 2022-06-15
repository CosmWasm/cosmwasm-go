package pkg

import (
	"errors"

	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
)

// CoinsContainMinAmount checks that coins have a target coin which amount is GTE to target.
func CoinsContainMinAmount(coins []stdTypes.Coin, coinExpected stdTypes.Coin) error {
	for _, coin := range coins {
		if coin.Denom != coinExpected.Denom {
			continue
		}

		if coin.Amount.LT(coinExpected.Amount) {
			break
		}

		return nil
	}

	return errors.New("expected coin amount (" + coinExpected.String() + "): not found")
}
