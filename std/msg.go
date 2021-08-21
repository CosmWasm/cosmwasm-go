package std

//------- Results / Msgs -------------

// HandleResult is the raw response from the handle call.
// This is mirrors Rust's ContractResult<HandleResponse>.
type HandleResult struct {
	Ok  *HandleResponse `json:",omitempty"`
	Err string          `json:"error,omitempty"`
}

// HandleResponse defines the return value on a successful handle
type HandleResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:",emptyslice"`
	// base64-encoded bytes to return as ABCI.Data field
	Data []byte `json:",omitempty"`
	// attributes for a log event to return over abci interface
	Attributes []EventAttribute `json:",emptyslice"`
}

// InitResult is the raw response from the handle call.
// This is mirrors Rust's ContractResult<InitResponse>.
type InitResult struct {
	Ok  *InitResponse `json:",omitempty"`
	Err string        `json:"error,omitempty"`
}

// InitResponse defines the return value on a successful handle
type InitResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:",emptyslice"`
	// attributes for a log event to return over abci interface
	Attributes []EventAttribute `json:",emptyslice"`
}

// MigrateResult is the raw response from the migrate call.
// This is mirrors Rust's ContractResult<MigrateResponse>.
type MigrateResult struct {
	Ok  *MigrateResponse `json:",omitempty"`
	Err string           `json:"error,omitempty"`
}

// MigrateResponse defines the return value on a successful handle
type MigrateResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:",emptyslice"`
	// base64-encoded bytes to return as ABCI.Data field
	Data []byte `json:",omitempty"`
	// attributes for a log event to return over abci interface
	Attributes []EventAttribute `json:",emptyslice"`
}

// EventAttribute
type EventAttribute struct {
	Key   string
	Value string
}

// CosmosMsg is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type CosmosMsg struct {
	Bank    *BankMsg    `json:",omitempty"`
	Custom  RawMessage  `json:",omitempty"`
	Staking *StakingMsg `json:",omitempty"`
	Wasm    *WasmMsg    `json:",omitempty"`
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
	Send []Coin `json:"send,emptyslice"`
}

type InstantiateMsg struct {
	// CodeID is the reference to the wasm byte code as used by the Cosmos-SDK
	CodeID uint64 `json:"code_id"`
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte `json:"msg"`
	// Send is an optional amount of coins this contract sends to the called contract
	Send []Coin `json:"send,emptyslice"`
}
