package src

import "github.com/cosmwasm/cosmwasm-go/std"

//all message type define here
type InitMsg struct {
	Count uint64 `json:"count"`
}

type HandleMsg struct {
	Increment Increment `json:"increment,omitempty"`
	Reset     Reset     `json:"reset,omitempty"`
}

type Increment struct {
	Delta uint64 `json:"delta"`
}

type Reset struct {
	Value uint64 `json:"value"`
}

type QueryMsg struct {
	Count std.EmptyStruct `json:"get_count"`
}

type CountResponse struct {
	Count uint64 `json:"count"`
}
