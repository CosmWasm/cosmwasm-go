package std

import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
	"strconv"
)

// Coin is a string representation of the sdk.Coin type (more portable than sdk.Int)
type Coin struct {
	Denom  string `json:"denom"`  // type, eg. "ATOM"
	Amount string `json:"amount"` // string encoing of decimal value, eg. "12.3456"
}

func NewCoin(amount uint64, denom string) Coin {
	return Coin{
		Denom:  denom,
		Amount: strconv.FormatUint(amount, 10),
	}
}

// Coins handles properly serializing empty amounts
type Coins []Coin

// MarshalJSON ensures that we get [] for empty arrays
func (c Coins) MarshalJSON() ([]byte, error) {
	if len(c) == 0 {
		return []byte("[]"), nil
	}
	var d []Coin = c
	return ezjson.Marshal(d)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (c *Coins) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var d []Coin
	if err := ezjson.Unmarshal(data, &d); err != nil {
		return err
	}
	*c = d
	return nil
}

// ============= MSG ===========
//------- Results / Msgs -------------

// CosmosResponse is the raw response from the init / handle calls
type CosmosResponse struct {
	Ok  *Result   `json:"Ok,omitempty"`
	Err *StdError `json:"Err,omitempty"`
}

// MarshalJSON returns m as the JSON encoding of m.
func (m CosmosResponse) MarshalJSON() ([]byte, error) {
	var str string
	if m.Ok != nil {
		str = `{"Ok":`+m.Ok.ToString()+"}"
	}else{
		str = `{"Err":`+m.Err.Error()+`}`
	}

	return []byte(str), nil
}

func CosmosResponseDefault() CosmosResponse {
	result := ResultDefault()
	return CosmosResponse{
		Ok:  &result,
		Err: nil,
	}
}

// Result defines the return value on a successful
type Result struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:"messages"`
	// base64-encoded bytes to return as ABCI.Data field
	Data string `json:"data"`
	// log message to return over abci interface
	Log []LogAttribute `json:"log"`
}

func (r Result) ToString() string {
	return `{"messages":[], "data":"`+r.Data+`","log":[]}`
}

func ResultDefault() Result {
	return Result{
		Messages: []CosmosMsg{},
		Data:     "",
		Log:      []LogAttribute{},
	}
}

// LogAttribute
type LogAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CosmosMsg is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type CosmosMsg struct {
	Bank    *BankMsg    `json:"bank,omitempty"`
	Custom  RawMessage  `json:"custom,omitempty"`
	Staking *StakingMsg `json:"staking,omitempty"`
	Wasm    *WasmMsg    `json:"wasm,omitempty"`
}

type BankMsg struct {
	Send *SendMsg `json:"send,omitempty"`
}

// SendMsg contains instructions for a Cosmos-SDK/SendMsg
// It has a fixed interface here and should be converted into the proper SDK format before dispatching
type SendMsg struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      Coins  `json:"amount"`
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
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

type StakingMsg struct {
	Delegate   *DelegateMsg   `json:"delegate,omitempty"`
	Undelegate *UndelegateMsg `json:"undelegate,omitempty"`
	Redelegate *RedelegateMsg `json:"redelegate,omitempty"`
	Withdraw   *WithdrawMsg   `json:"withdraw,omitempty"`
}

type DelegateMsg struct {
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

type UndelegateMsg struct {
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

type RedelegateMsg struct {
	SrcValidator string `json:"src_validator"`
	DstValidator string `json:"dst_validator"`
	Amount       Coin   `json:"amount"`
}

type WithdrawMsg struct {
	Validator string `json:"validator"`
	// this is optional
	Recipient string `json:"recipient,omitempty"`
}

type WasmMsg struct {
	Execute     *ExecuteMsg     `json:"execute,omitempty"`
	Instantiate *InstantiateMsg `json:"instantiate,omitempty"`
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
	ContractAddr string `json:"contract_addr"`
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte `json:"msg"`
	// Send is an optional amount of coins this contract sends to the called contract
	Send Coins `json:"send"`
}

type InstantiateMsg struct {
	// CodeID is the reference to the wasm byte code as used by the Cosmos-SDK
	CodeID uint64 `json:"code_id"`
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte `json:"msg"`
	// Send is an optional amount of coins this contract sends to the called contract
	Send Coins `json:"send"`
}
