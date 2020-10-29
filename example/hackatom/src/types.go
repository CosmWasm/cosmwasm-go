package src

// this is what we store
type State struct {
	Count uint64 `json:"count"`
	Owner []byte `json:"owner"`
}

var StateKey = []byte("State")

//all message type define here
type InitMsg struct {
	Count uint64 `json:"count"`
}

type Handler struct {
	Increment              Increment               `json:"increment,omitempty"`
	Reset          Reset           `json:"reset,omitempty"`
}

type Increment struct {
	Delta uint64 `json:"delta"`
}

type Reset struct {
	Value uint64 `json:"value"`
}

type Querier struct {
	Count Count `json:"get_count"`
}

// how to better handle empty values for enum? it must be different than default/zero
type Count struct {
	A string `json:"a"`
}

type CountResponse struct {
	Count uint64 `json:"count"`
}
