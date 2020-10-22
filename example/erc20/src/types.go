package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
)

//all message type define here
type InitMsg struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimal     uint64 `json:"decimal"`
	TotalSupply uint64 `json:"total_supply"`
}

// returns an error if invalid
func (i InitMsg) Validate() *std.CosmosResponseError {
	if len(i.Name) < 2 {
		return std.GenerateError(std.GenericError, "Name must be at least 2 characters", "")
	}
	if len(i.Symbol) < 2 || len(i.Symbol) > 6 {
		return std.GenerateError(std.GenericError, "Symbol must be 2-6 characters", "")
	}
	if i.TotalSupply == 0 {
		return std.GenerateError(std.GenericError, "Total Supply must be greater than 0", "")
	}
	return nil
}

type Transfer struct {
	ToAddr string `json:"to"`
	Value  uint64 `json:"value"`
}

type TransferFrom struct {
	FromAddr string `json:"from"`
	ToAddr   string `json:"to"`
	Value    uint64 `json:"value"`
}

type Approve struct {
	Spender string `json:"spender"`
	Value   uint64 `json:"value"`
}

type TransferOwner struct {
	NewOwner string `json:"new_owner"`
}

type AcceptTransferredOwner struct {
	AcceptedAddr string `json:"accepted_address"`
}

type Handler struct {
	Transfer              Transfer               `json:"transfer,omitempty"`
	TransferFrom          TransferFrom           `json:"transfer_from,omitempty"`
	Approve               Approve                `json:"approve,omitempty"`
	TransferOwner         TransferOwner          `json:"transfer_owner,omitempty"`
	AcceptedTransferOwner AcceptTransferredOwner `json:"accept_transfer_owner,omitempty"`
}

type Balance struct {
	Address []byte `json:"address"`
}

type BalanceResponse struct {
	Value uint64 `json:"value"`
}

type Querier struct {
	Balance Balance `json:"balance"`
}
