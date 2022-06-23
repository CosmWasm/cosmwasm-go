package types

import (
	"errors"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/pkg"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
)

// MsgSudo is handled by the Sudo entrypoint.
type MsgSudo struct {
	ChangeNewVotingCost *ChangeCostRequest `json:",omitempty"`
	ChangeVoteCost      *ChangeCostRequest `json:",omitempty"`
}

// ChangeCostRequest defines MsgSudo.ChangeNewVotingCost and MsgSudo.ChangeVoteCost request.
type ChangeCostRequest struct {
	NewCost stdTypes.Coin
}

// Validate performs object fields validation.
func (r ChangeCostRequest) Validate() error {
	if err := pkg.ValidateDenom(r.NewCost.Denom); err != nil {
		return errors.New("newCost: " + err.Error())
	}

	return nil
}
