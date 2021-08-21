package std

import (
	"strconv"
)

// Coin is a string representation of the sdk.Coin type (more portable than sdk.Int)
type Coin struct {
	Denom  string // type, eg. "ATOM"
	Amount string // string encoing of decimal value, eg. "12.3456"
}

func NewCoin(amount uint64, denom string) Coin {
	return Coin{
		Denom:  denom,
		Amount: strconv.FormatUint(amount, 10),
	}
}

func (c Coin) IsEmpty() bool {
	return c.Denom == "" && c.Amount == ""
}

// We preallocate empty elements at the end for parsing.
// This will remove the ones that were not filled
func TrimCoins(parsed []Coin) []Coin {
	i := 0
	for !parsed[i].IsEmpty() {
		i++
	}
	return parsed[:i]
}

func NewCoins(amount uint64, denom string) []Coin {
	return []Coin{NewCoin(amount, denom)}
}

// ============= MSG ===========
//------- Results / Msgs -------------

// InitResponse defines the return value on a successful handle
type InitResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:",emptyslice"`
	// log message to return over abci interface
	Attributes []Attribute `json:",emptyslice"`
}

type InitResultOk struct {
	Ok InitResponse `json:"ok"`
}

func InitResultOkOkDefault() *InitResultOk {
	return &InitResultOk{
		Ok: InitResponse{
			Messages:   []CosmosMsg{},
			Attributes: []Attribute{},
		},
	}
}

type CosmosResponseError struct {
	Err StdError
}

type OptionBinary struct {
	None string `json:",omitempty"`
	Some []byte `json:",omitempty"`
}

// HandleResponse defines the return value on a successful handle
type HandleResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:",emptyslice"`
	// base64-encoded bytes to return as ABCI.Data field
	Data string `json:",omitempty"`
	// log message to return over abci interface
	Attributes []Attribute `json:",emptyslice"`
}

type HandleResultOk struct {
	Ok HandleResponse `json:"ok"`
}

func HandleResultOkDefault() *HandleResultOk {
	return &HandleResultOk{
		Ok: HandleResponse{
			Messages:   []CosmosMsg{},
			Attributes: []Attribute{},
			Data:       "",
		},
	}
}

// MigrateResponse defines the return value on a successful handle
type MigrateResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:",emptyslice"`
	// base64-encoded bytes to return as ABCI.Data field
	Data string `json:",omitempty"`
	// log message to return over abci interface
	Attributes []Attribute `json:",emptyslice"`
}

type MigrateResultOk struct {
	Ok MigrateResponse `json:"ok"`
}

func MigrateResultOkDefault() *MigrateResultOk {
	return &MigrateResultOk{
		Ok: MigrateResponse{
			Messages:   []CosmosMsg{},
			Attributes: []Attribute{},
			Data:       "",
		},
	}
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CosmosMsg is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type CosmosMsg struct {
	Bank    *BankMsg    `json:"bank,omitempty"`
	Custom  *RawMessage `json:"custom,omitempty"`
	Staking *StakingMsg `json:"staking,omitempty"`
	Wasm    *WasmMsg    `json:"wasm,omitempty"`
}

type BankMsg struct {
	Send *SendMsg `json:",omitempty"`
}

// SendMsg contains instructions for a Cosmos-SDK/SendMsg
// It has a fixed interface here and should be converted into the proper SDK format before dispatching
type SendMsg struct {
	FromAddress string
	ToAddress   string
	Amount      []Coin `json:",emptyslice"`
}

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return NewError("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

type StakingMsg struct {
	Delegate   *DelegateMsg   `json:",omitempty"`
	Undelegate *UndelegateMsg `json:",omitempty"`
	Redelegate *RedelegateMsg `json:",omitempty"`
	Withdraw   *WithdrawMsg   `json:",omitempty"`
}

type DelegateMsg struct {
	Validator string
	Amount    Coin
}

type UndelegateMsg struct {
	Validator string
	Amount    Coin
}

type RedelegateMsg struct {
	SrcValidator string
	DstValidator string
	Amount       Coin
}

type WithdrawMsg struct {
	Validator string
	// this is optional
	Recipient string `json:",omitempty"`
}

type WasmMsg struct {
	Execute     *ExecuteMsg     `json:",omitempty"`
	Instantiate *InstantiateMsg `json:",omitempty"`
}

// ExecuteMsg is used to call another defined contract on this chain.
// The calling contract requires the callee to be defined beforehand,
// and the address should have been defined in initialization.
// And we assume the developer tested the ABIs and coded them together.
//
// Since a contract is immutable once it is deployed, we don't need to transform this.
// If it was properly coded and worked once, it will continue to work throughout upgrades.
type ExecuteMsg struct {
	// ContractAddr is the sdk.AccAddress of the contract, which uniquely defines
	// the contract ID and instance ID. The sdk module should maintain a reverse lookup table.
	ContractAddr string
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte
	// Send is an optional amount of coins this contract sends to the called contract
	Send []Coin `json:",emptyslice"`
}

type InstantiateMsg struct {
	// CodeID is the reference to the wasm byte code as used by the Cosmos-SDK
	CodeID uint64
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte
	// Send is an optional amount of coins this contract sends to the called contract
	Send []Coin `json:",emptyslice"`
}
