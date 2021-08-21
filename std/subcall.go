package std

type replyOn int

const (
	ReplyAlways  = "always"
	ReplySuccess = "success"
	ReplyError   = "error"
	ReplyNever   = "never"
)

// SubMsg wraps a CosmosMsg with some metadata for handling replies (ID) and optionally
// limiting the gas usage (GasLimit)
type SubMsg struct {
	ID       uint64    `json:"id"`
	Msg      CosmosMsg `json:"msg"`
	GasLimit *uint64   `json:"gas_limit,omitempty"`
	ReplyOn  string    `json:"reply_on"`
}

type Reply struct {
	ID     uint64        `json:"id"`
	Result SubcallResult `json:"result"`
}

// SubcallResult is the raw response we return from the sdk -> reply after executing a SubMsg.
// This is mirrors Rust's ContractResult<SubcallResponse>.
type SubcallResult struct {
	Ok  *SubcallResponse `json:"ok,omitempty"`
	Err string           `json:"error,omitempty"`
}

type SubcallResponse struct {
	Events []Event `json:"events,emptyslice"`
	Data   []byte  `json:"data,omitempty"`
}

type Event struct {
	Type       string           `json:"type"`
	Attributes []EventAttribute `json:"attributes,emptyslice"`
}

func NewSubMsg(msg CosmosMsg) SubMsg {
	return SubMsg{
		ID:      0,
		Msg:     msg,
		ReplyOn: ReplyNever,
	}
}
