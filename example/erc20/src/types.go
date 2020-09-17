package src

//all message type define here
type InitMsg struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimal     uint64 `json:"decimal"`
	TotalSupply uint64 `json:"total_supply"`
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
