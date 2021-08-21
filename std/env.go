package std

//---------- Env ---------

// Env defines the state of the blockchain environment this contract is
// running in. This must contain only trusted data - nothing from the Tx itself
// that has not been verfied (like Signer).
//
// Env are json encoded to a byte slice before passing to the wasm contract.
type Env struct {
	Block    BlockInfo
	Contract ContractInfo
}

type BlockInfo struct {
	// block height this transaction is executed
	Height uint64
	// time in seconds since unix epoch (since CosmWasm 0.3)
	Time uint64
	// Nanoseconds of the block time (since CosmWasm 0.11)
	TimeNanos uint64
	ChainID   string
}

type MessageInfo struct {
	// bech32 encoding of sdk.AccAddress executing the contract
	Sender string
	// amount of funds send to the contract along with this message
	SentFunds []Coin
}

type ContractInfo struct {
	// bech32 encoding of sdk.AccAddress of the contract, to be used when sending messages
	Address string
}
