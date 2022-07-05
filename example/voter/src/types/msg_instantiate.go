package types

import (
	"errors"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/pkg"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
)

// MsgInstantiate is handled by the Instantiate entrypoint.
type MsgInstantiate struct {
	// Params are the contract parameters.
	Params state.Params
}

// Validate performs object fields validation.
func (m MsgInstantiate) Validate(info stdTypes.MessageInfo) error {
	if m.Params.OwnerAddr != info.Sender {
		return errors.New("params.ownerAddr: must EQ to senderAddr")
	}

	if err := pkg.ValidateDenom(m.Params.NewVotingCost.Denom); err != nil {
		return errors.New("params.newVotingCost: " + err.Error())
	}

	if err := pkg.ValidateDenom(m.Params.VoteCost.Denom); err != nil {
		return errors.New("params.voteCost: " + err.Error())
	}

	if m.Params.IBCSendTimeout == 0 {
		return errors.New("params.ibcSendTimeout: must be GT 0")
	}

	return nil
}
