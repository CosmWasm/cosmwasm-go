package src

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
	Count Count `json:"get_count"`
}

// how to better handle empty values for enum? it must be different than default/zero
type Count struct {
	A string `json:"a"`
}

type CountResponse struct {
	Count int64 `json:"count"`
}
