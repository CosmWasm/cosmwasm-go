package std

import (
	"encoding/base64"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// ------- query detail types ---------
type QueryResponse struct {
	// this must be base64 encoded
	Ok string `json:"ok,omitempty,rust_option"`
	// TODO: what is this format actually?
	Error string `json:"error,omitempty"`
}

// This is a 2-level result
type QuerierResult struct {
	Ok QueryResponse `json:"ok,omitempty"`
	// TODO: what is this format actually?
	Error string `json:"error,omitempty"`
}

func BuildQueryResponse(msg string) *QueryResponse {
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	return &QueryResponse{Ok: encoded}
}

func BuildQueryResponseBinary(msg []byte) *QueryResponse {
	encoded := base64.StdEncoding.EncodeToString(msg)
	return &QueryResponse{Ok: encoded}
}

func (q QueryResponse) Data() ([]byte, error) {
	if q.Error != "" {
		return nil, errors.New(q.Error)
	}
	return base64.StdEncoding.DecodeString(q.Ok)
}

//
//// Query will handle most of the marshalling/unmarshalling. You need to parse the final result
//func Query(querier Querier, request QueryRequest) ([]byte, error) {
//	raw, err := ezjson.Marshal(request)
//	if err != nil {
//		return nil, err
//	}
//	querier.RawQuery(raw)
//}

// QueryRequest is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type QueryRequest struct {
	Bank    BankQuery    `json:"bank,omitempty"`
	Custom  RawMessage   `json:"custom,omitempty"`
	Staking StakingQuery `json:"staking,omitempty"`
	Wasm    WasmQuery    `json:"wasm,omitempty"`
}

type BankQuery struct {
	Balance     BalanceQuery     `json:"balance,omitempty"`
	AllBalances AllBalancesQuery `json:"all_balances,omitempty"`
}

func (b BankQuery) IsEmpty() bool {
	return b.Balance.Address == "" && b.AllBalances.Address == ""
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
	Amount Coins `json:"amount"`
}

type StakingQuery struct {
	Validators     ezjson.EmptyStruct  `json:"validators,omitempty,opt_seen"`
	AllDelegations AllDelegationsQuery `json:"all_delegations,omitempty"`
	Delegation     DelegationQuery     `json:"delegation,omitempty"`
	BondedDenom    ezjson.EmptyStruct  `json:"bonded_denom,omitempty,opt_seen"`
}

func (s StakingQuery) IsEmpty() bool {
	return !s.Validators.WasSet() &&
		s.AllDelegations.Delegator == "" &&
		s.Delegation.Delegator == "" &&
		!s.BondedDenom.WasSet()
}

// ValidatorsResponse is the expected response to ValidatorsQuery
type ValidatorsResponse struct {
	Validators Validators `json:"validators"`
}

// TODO: Validators must JSON encode empty array as []
type Validators []Validator

// MarshalJSON ensures that we get [] for empty arrays
func (v Validators) MarshalJSON() ([]byte, error) {
	if len(v) == 0 {
		return []byte("[]"), nil
	}
	var raw []Validator = v
	return ezjson.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (v *Validators) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Validator
	if err := ezjson.Unmarshal(data, &raw); err != nil {
		return err
	}
	*v = raw
	return nil
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
	Delegations Delegations `json:"delegations"`
}

type Delegations []Delegation

// MarshalJSON ensures that we get [] for empty arrays
func (d Delegations) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	var raw []Delegation = d
	return ezjson.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (d *Delegations) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Delegation
	if err := ezjson.Unmarshal(data, &raw); err != nil {
		return err
	}
	*d = raw
	return nil
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
	AccumulatedRewards []Coin `json:"accumulated_rewards"`
	CanRedelegate      Coin   `json:"can_redelegate"`
}

type BondedDenomResponse struct {
	Denom string `json:"denom"`
}

type WasmQuery struct {
	Smart SmartQuery `json:"smart,omitempty"`
	Raw   RawQuery   `json:"raw,omitempty"`
}

func (w WasmQuery) IsEmpty() bool {
	return w.Smart.ContractAddr == "" &&
		w.Raw.ContractAddr == ""
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
