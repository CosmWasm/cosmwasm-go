package std

import (
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

// Deps contains the dependencies passed to a contract's mutating entrypoints.
type Deps struct {
	// Storage provides access to the data persistence layer at write and read level.
	Storage Storage
	// Api provides access to common utilities such as address
	// parsing and verification.
	Api Api
	// Querier is used to query information from other contracts.
	Querier Querier
}

type Order uint32

const (
	Ascending  Order = 1
	Descending Order = 2
)

// ReadonlyStorage defines the behaviour of a KV with only read capabilities.
type ReadonlyStorage interface {
	// Get gets the value of the provided key. If value is nil then the key does not exist.
	Get(key []byte) (value []byte)
	// Range ranges from start to end byte prefixes with the provided Order flag.
	Range(start, end []byte, order Order) (iterator Iterator)
}

// Storage defines the behaviour of a KV with read and write capabilities.
type Storage interface {
	ReadonlyStorage

	// Set sets the key and value.
	Set(key, value []byte)
	// Remove removes the value from the db.
	Remove(key []byte)
}

type Iterator interface {
	Next() (key, value []byte, err error)
}

type Api interface {
	CanonicalAddress(human types.HumanAddress) (types.CanonicalAddress, error)
	HumanAddress(canonical types.CanonicalAddress) (types.HumanAddress, error)
	ValidateAddress(human types.HumanAddress) error
	Debug(msg string)
}

type Querier interface {
	RawQuery(request []byte) ([]byte, error)
}

type QuerierWrapper struct {
	Querier
}

type JSONType interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

func (q QuerierWrapper) doQuery(query types.QueryRequest, result JSONType) error {
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

func (q QuerierWrapper) QueryAllBalances(addr string) ([]types.Coin, error) {
	query := types.QueryRequest{
		Bank: &types.BankQuery{
			AllBalances: &types.AllBalancesQuery{
				Address: addr,
			},
		},
	}
	qres := types.AllBalancesResponse{}
	err := q.doQuery(query, &qres)
	if err != nil {
		return nil, err
	}
	return qres.Amount, nil
}

func (q QuerierWrapper) QueryBalance(addr string, denom string) (types.Coin, error) {
	query := types.QueryRequest{
		Bank: &types.BankQuery{
			Balance: &types.BalanceQuery{
				Address: addr,
				Denom:   denom,
			},
		},
	}
	qres := types.BalanceResponse{}
	err := q.doQuery(query, &qres)
	return qres.Amount, err
}
