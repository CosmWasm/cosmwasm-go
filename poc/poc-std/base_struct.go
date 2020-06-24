package poc_std

import "math/big"

type BlockInfo struct {
	Height uint64
	Time uint64
	Chain_id string
}

type Coin struct {
	denom string
	amount big.Int
}

type MessageInfo struct {
	Sender string
	Sent_funds []Coin
}

type ContractInfo struct {
	address string
}

type Env struct {
	Block BlockInfo
	Message MessageInfo
	Contract ContractInfo
}
