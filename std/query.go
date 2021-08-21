package std

import (
	"encoding/base64"
	"errors"
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

type QuerierWrapper struct {
	Querier
}

type JSONType interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

func (q QuerierWrapper) doQuery(query QueryRequest, result JSONType) error {
	binQuery, err := query.MarshalJSON()
	if err != nil {
		return err
	}
	data, err := q.Querier.RawQuery(binQuery)
	if err != nil {
		return err
	}
	return result.UnmarshalJSON(data)
}

func (q QuerierWrapper) QueryAllBalances(addr string) ([]Coin, error) {
	query := QueryRequest{
		Bank: &BankQuery{
			AllBalances: &AllBalancesQuery{
				Address: addr,
			},
		},
	}
	qres := AllBalancesResponse{}
	err := q.doQuery(query, &qres)
	if err != nil {
		return nil, err
	}
	return qres.Amount, nil
}

func (q QuerierWrapper) QueryBalance(addr string, denom string) (Coin, error) {
	query := QueryRequest{
		Bank: &BankQuery{
			Balance: &BalanceQuery{
				Address: addr,
				Denom:   denom,
			},
		},
	}
	qres := BalanceResponse{}
	err := q.doQuery(query, &qres)
	return qres.Amount, err
}

// QueryRequest is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type QueryRequest struct {
	Bank    *BankQuery    `json:"bank,omitempty"`
	Custom  *RawMessage   `json:"custom,omitempty"`
	Staking *StakingQuery `json:"staking,omitempty"`
	Wasm    *WasmQuery    `json:"wasm,omitempty"`
}

type BankQuery struct {
	Balance     *BalanceQuery     `json:"balance,omitempty"`
	AllBalances *AllBalancesQuery `json:"all_balances,omitempty"`
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
	Amount []Coin `json:"amount,emptyslice"`
}

type StakingQuery struct {
	Validators     *struct{}            `json:",omitempty"`
	AllDelegations *AllDelegationsQuery `json:",omitempty"`
	Delegation     *DelegationQuery     `json:",omitempty"`
	BondedDenom    *struct{}            `json:",omitempty"`
}

// ValidatorsResponse is the expected response to ValidatorsQuery
type ValidatorsResponse struct {
	Validators []Validator `json:",emptyslice"`
}
type Validator struct {
	Address string
	// decimal string, eg "0.02"
	Commission string
	// decimal string, eg "0.02"
	MaxCommission string
	// decimal string, eg "0.02"
	MaxChangeRate string
}

type AllDelegationsQuery struct {
	Delegator string
}

type DelegationQuery struct {
	Delegator string
	Validator string
}

// AllDelegationsResponse is the expected response to AllDelegationsQuery
type AllDelegationsResponse struct {
	Delegations []Delegation `json:",emptyslice"`
}
type Delegation struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

// DelegationResponse is the expected response to DelegationsQuery
type DelegationResponse struct {
	Delegation *FullDelegation `json:",omitempty"`
}

type FullDelegation struct {
	Delegator          string
	Validator          string
	Amount             Coin
	AccumulatedRewards []Coin `json:",emptyslice"`
	CanRedelegate      Coin
}

type BondedDenomResponse struct {
	Denom string
}

type WasmQuery struct {
	Smart *SmartQuery `json:",omitempty"`
	Raw   *RawQuery   `json:",omitempty"`
}

// SmartQuery respone is raw bytes ([]byte)
type SmartQuery struct {
	ContractAddr string
	Msg          []byte
}

// RawQuery response is raw bytes ([]byte)
type RawQuery struct {
	ContractAddr string
	Key          []byte
}
