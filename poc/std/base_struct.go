package std

type BlockInfo struct {
	Height uint64
	Time uint64
	Chain_id string
}

type Coin struct {
	denom string
	amount uint64
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

type MemRegion struct {
	Offset 		uint32
	Capacity 	uint32
	Length		uint32
}

const REGION_HEAD_SIZE uint32 = 12
