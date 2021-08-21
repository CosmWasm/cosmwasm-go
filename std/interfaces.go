package std

// =========== Deps --> context =======
type Deps struct {
	Storage Storage
	Api     Api
	Querier Querier
}

type Order uint32

const (
	Ascending  Order = 1
	Descending Order = 2
)

type ReadonlyStorage interface {
	Get(key []byte) (value []byte, err error)
	Range(start, end []byte, order Order) (Iterator, error)
}

type Storage interface {
	ReadonlyStorage

	Set(key, value []byte) error
	Remove(key []byte) error
}

type Iterator interface {
	Next() (key, value []byte, err error)
}

type Api interface {
	CanonicalAddress(human string) (CanonicalAddr, error)
	HumanAddress(canonical CanonicalAddr) (string, error)
	ValidateAddress(human string) error
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
