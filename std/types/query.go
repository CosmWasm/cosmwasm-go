package types

// ------- query detail types ---------
// QueryResponse is the Go counterpart of `ContractResult<Binary>`.
// The JSON annotations are used for deserializing directly. There is a custom serializer below.
// type QueryResponse queryResponseImpl

// type queryResponseImpl struct {
type QueryResponse struct {
	Ok  []byte `json:"ok,omitempty"`
	Err string `json:"error,omitempty"`
}

// TODO: after tinygo gen, try to ensure we handle this properly (all empty as ok)
// Or maybe we only need that for Rust contracts to parse fine??

// // A custom serializer that allows us to map QueryResponse instances to the Rust
// // enum `ContractResult<Binary>`
// func (q QueryResponse) MarshalJSON() ([]byte, error) {
// 	// In case both Ok and Err are empty, this is interpreted and seralized
// 	// as an Ok case with no data because errors must not be empty.
// 	if len(q.Ok) == 0 && len(q.Err) == 0 {
// 		return []byte(`{"ok":""}`), nil
// 	}
// 	return queryResponseImpl(q).MarshalJSON()
// }

// This is a 2-level result
type QuerierResult struct {
	Ok  *QueryResponse `json:"ok,omitempty"`
	Err *SystemError   `json:"error,omitempty"`
}

// TODO: I think this can be removed (we may need to update calling code)
func BuildQueryResponse(msg string) *QueryResponse {
	return &QueryResponse{Ok: []byte(msg)}
}

func BuildQueryResponseBinary(msg []byte) *QueryResponse {
	return &QueryResponse{Ok: msg}
}

func (q QueryResponse) Data() ([]byte, error) {
	if q.Err != "" {
		return nil, GenericError(q.Err)
	}
	return q.Ok, nil
}

// QueryRequest is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type QueryRequest struct {
	Bank     *BankQuery     `json:"bank,omitempty"`
	Custom   RawMessage     `json:"custom,omitempty"`
	IBC      *IBCQuery      `json:"ibc,omitempty"`
	Staking  *StakingQuery  `json:"staking,omitempty"`
	Stargate *StargateQuery `json:"stargate,omitempty"`
	Wasm     *WasmQuery     `json:"wasm,omitempty"`
}

type BankQuery struct {
	Balance     *BalanceQuery     `json:"balance,omitempty"`
	AllBalances *AllBalancesQuery `json:"all_balances,omitempty"`
}

type BalanceQuery struct {
	Address string `json:"address"`
	Denom   string `json:"denom"`
}

// BalanceResponse is the expected response to BalanceQuery
type BalanceResponse struct {
	Amount Coin `json:"amount"`
}

type AllBalancesQuery struct {
	Address string `json:"address"`
}

// AllBalancesResponse is the expected response to AllBalancesQuery
type AllBalancesResponse struct {
	Amount []Coin `json:"amount,emptyslice"`
}

type StakingQuery struct {
	AllValidators  *AllValidatorsQuery  `json:"all_validators,omitempty"`
	Validator      *ValidatorQuery      `json:"validator,omitempty"`
	AllDelegations *AllDelegationsQuery `json:"all_delegations,omitempty"`
	Delegation     *DelegationQuery     `json:"delegation,omitempty"`
	BondedDenom    *struct{}            `json:"bonded_denom,omitempty"`
}

type AllValidatorsQuery struct{}

// AllValidatorsResponse is the expected response to AllValidatorsQuery
type AllValidatorsResponse struct {
	Validators []Validator `json:"validators,emptyslice"`
}
type ValidatorQuery struct {
	/// Address is the validator's address (e.g. cosmosvaloper1...)
	Address string `json:"address"`
}

// ValidatorResponse is the expected response to ValidatorQuery
type ValidatorResponse struct {
	Validator *Validator `json:"validator"` // serializes to `null` when unset which matches Rust's Option::None serialization
}

type Validator struct {
	Address string `json:"address"`
	// decimal string, eg "0.02"
	Commission string `json:"commission"`
	// decimal string, eg "0.02"
	MaxCommission string `json:"max_commission"`
	// decimal string, eg "0.02"
	MaxChangeRate string `json:"max_change_rate"`
}

type AllDelegationsQuery struct {
	Delegator string `json:"delegator"`
}

type DelegationQuery struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
}

// AllDelegationsResponse is the expected response to AllDelegationsQuery
type AllDelegationsResponse struct {
	Delegations []Delegation `json:"delegations,emptyslice"`
}
type Delegation struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

// DelegationResponse is the expected response to DelegationsQuery
type DelegationResponse struct {
	Delegation *FullDelegation `json:"delegation,omitempty"`
}

type FullDelegation struct {
	Delegator          string `json:"delegator"`
	Validator          string `json:"validator"`
	Amount             Coin   `json:"amount"`
	AccumulatedRewards []Coin `json:"accumulated_rewards,emptyslice"`
	CanRedelegate      Coin   `json:"can_redelegate"`
}

type BondedDenomResponse struct {
	Denom string `json:"denom"`
}

// A Stargate query encoded the same way as abci_query, with path and protobuf encoded Data.
// The format is defined in [ADR-21](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-021-protobuf-query-encoding.md)
// The response is also protobuf encoded. The caller is responsible for compiling the proper protobuf definitions
type StargateQuery struct {
	// this is the fully qualified service path used for routing,
	// eg. custom/cosmos_sdk.x.bank.v1.Query/QueryBalance
	Path string `json:"path"`
	// this is the expected protobuf message type (not any), binary encoded
	Data []byte `json:"data"`
}

// This is the protobuf response, binary encoded.
// The caller is responsible for knowing how to parse.
type StargateResponse struct {
	Response []byte `json:"response"`
}

type WasmQuery struct {
	Smart        *SmartQuery        `json:"smart,omitempty"`
	Raw          *RawQuery          `json:"raw,omitempty"`
	ContractInfo *ContractInfoQuery `json:"contract_info,omitempty"`
}

// SmartQuery respone is raw bytes ([]byte)
type SmartQuery struct {
	ContractAddr string `json:"contract_addr"`
	Msg          []byte `json:"msg"`
}

// RawQuery response is raw bytes ([]byte)
type RawQuery struct {
	ContractAddr string `json:"contract_addr"`
	Key          []byte `json:"key"`
}

type ContractInfoQuery struct {
	// Bech32 encoded sdk.AccAddress of the contract
	ContractAddr string `json:"contract_addr"`
}

type ContractInfoResponse struct {
	CodeID  uint64 `json:"code_id"`
	Creator string `json:"creator"`
	// Set to the admin who can migrate contract, if any
	Admin  string `json:"admin,omit_empty"`
	Pinned bool   `json:"pinned"`
	// Set if the contract is IBC enabled
	IBCPort string `json:"ibc_port,omit_empty"`
}
