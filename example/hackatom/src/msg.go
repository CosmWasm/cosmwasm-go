package src

import "github.com/cosmwasm/cosmwasm-go/std/ezjson"

//all message type define here
type InitMsg struct {
	Count uint64 `json:"count"`
}

type HandleMsg struct {
	Increment ezjson.EmptyStruct `json:"increment,opt_seen"`
	Reset     Reset              `json:"reset"`
}

type Increment struct {
	Delta uint64 `json:"delta"`
}

type Reset struct {
	Value uint64 `json:"value"`
}

type QueryMsg struct {
	Count ezjson.EmptyStruct `json:"get_count,opt_seen"`
}

type CountResponse struct {
	Count uint64 `json:"count"`
}
