package std

//---------- Env ---------

// Env defines the state of the blockchain environment this contract is
// running in. This must contain only trusted data - nothing from the Tx itself
// that has not been verfied (like Signer).
//
// Env are json encoded to a byte slice before passing to the wasm contract.
type Env struct {
	Block    BlockInfo    `block`
	Message  MessageInfo  `message`
	Contract ContractInfo `contract`
}

type BlockInfo struct {
	// block height this transaction is executed
	Height uint64 `height`
	// time in seconds since unix epoch - since cosmwasm 0.3
	Time    uint64 `time`
	ChainID string `chain_id`
}

type MessageInfo struct {
	// binary encoding of sdk.AccAddress executing the contract
	Sender []byte `sender`
	// amount of funds send to the contract along with this message
	SentFunds []Coin `sent_funds`
}

type ContractInfo struct {
	// binary encoding of sdk.AccAddress of the contract, to be used when sending messages
	Address []byte `address`
}
