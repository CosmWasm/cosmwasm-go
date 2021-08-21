// +build !cosmwasm

package std

import (
	"fmt"

	dbm "github.com/tendermint/tm-db"
)

func MockDeps(funds []Coin) *Deps {
	return &Deps{
		Storage: NewMockStorage(),
		Api:     MockApi{},
		Querier: NewMockQuerier(funds),
	}
}

const MOCK_CONTRACT_ADDR = "test-contract"

func MockEnv() Env {
	return Env{
		Block: BlockInfo{
			Height:    12_345,
			Time:      1_571_797_419,
			TimeNanos: 404_808_777,
			ChainID:   "cosmos-testnet-14002",
		},
		Contract: ContractInfo{
			Address: MOCK_CONTRACT_ADDR,
		},
	}
}

func MockInfo(sender string, funds []Coin) MessageInfo {
	return MessageInfo{
		Sender:    sender,
		SentFunds: funds,
	}
}

var (
	_ Iterator = (*MockIterator)(nil)
)

type MockIterator struct {
	Iter dbm.Iterator
}

func newMockIterator(iter dbm.Iterator) MockIterator {
	return MockIterator{
		Iter: iter,
	}
}

func (iter MockIterator) Next() (key, value []byte, err error) {
	if !iter.Iter.Valid() {
		iter.Iter.Close()
		return key, value, NewError("the end of iterator")
	}
	key, value = iter.Iter.Key(), iter.Iter.Value()
	iter.Iter.Next()
	return
}

type MockStorage struct {
	storage dbm.DB
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		storage: dbm.NewMemDB(),
	}
}

var (
	_ ReadonlyStorage = (*MockStorage)(nil)
	_ Storage         = (*MockStorage)(nil)
)

func (s *MockStorage) Get(key []byte) ([]byte, error) {
	return s.storage.Get(key)
}

func (s *MockStorage) Range(start, end []byte, order Order) (iter Iterator, err error) {
	var iterator dbm.Iterator
	switch order {
	case Ascending:
		iterator, err = s.storage.Iterator(start, end)
		iter = newMockIterator(iterator)
	case Descending:
		iterator, err = s.storage.ReverseIterator(start, end)
		iter = newMockIterator(iterator)
	default:
		err = NewError("failed. unexpected Order")
	}
	return
}

func (s *MockStorage) Set(key, value []byte) error {
	return s.storage.Set(key, value)
}

func (s *MockStorage) Remove(key []byte) error {
	return s.storage.Delete(key)
}

type CanonicalAddr []byte

const canonicalLength = 32

// ensure Api interface compliance at compile time
var (
	_ Api = (*MockApi)(nil)
)

type MockApi struct{}

func (api MockApi) CanonicalAddress(human string) (CanonicalAddr, error) {
	if len(human) > canonicalLength {
		return nil, NewError("failed. human encoding too long")
	}

	return []byte(human), nil
}

func (api MockApi) HumanAddress(canonical CanonicalAddr) (string, error) {
	if len(canonical) != canonicalLength {
		return "", NewError("failed. wrong canonical address length")
	}

	cutIndex := canonicalLength
	for i, v := range canonical {
		if v == 0 {
			cutIndex = i
			break
		}
	}

	return string(canonical[:cutIndex]), nil
}

func (api MockApi) Debug(msg string) {
	fmt.Println("DEBUG: " + msg)
}

// ====== Querier ======

// ensure Api interface compliance at compile time
var (
	_ Querier = (*MockQuerier)(nil)
)

type MockQuerier struct {
	Balances map[string][]Coin
}

func NewMockQuerier(funds []Coin) *MockQuerier {
	q := MockQuerier{
		Balances: make(map[string][]Coin),
	}
	if len(funds) > 0 {
		q.SetBalance(MOCK_CONTRACT_ADDR, funds)
	}
	return &q
}

func (q *MockQuerier) RawQuery(raw []byte) ([]byte, error) {
	var request QueryRequest
	err := request.UnmarshalJSON(raw)
	if err != nil {
		return nil, err
	}
	res, err := q.HandleQuery(request)
	if err != nil {
		return nil, err
	}
	return res.MarshalJSON()
}

func (q *MockQuerier) HandleQuery(request QueryRequest) (JSONType, error) {
	switch {
	case request.Bank != nil:
		return q.HandleBank(request.Bank)
	case request.Staking != nil:
		return nil, NewError("Staking queries not implemented")
	case request.Wasm != nil:
		return nil, NewError("Wasm queries not implemented")
	case request.Custom != nil:
		return nil, NewError("Custom queries not implemented")
	default:
		return nil, NewError("Unknown QueryRequest variant")
	}
}

func (q *MockQuerier) HandleBank(request *BankQuery) (JSONType, error) {
	switch {
	case request.Balance != nil:
		balances := q.GetBalance(request.Balance.Address)
		coin := Coin{Denom: request.Balance.Denom, Amount: "0"}
		for _, c := range balances {
			if c.Denom == coin.Denom {
				coin.Amount = c.Amount
				break
			}
		}
		return &BalanceResponse{Amount: coin}, nil
	case request.AllBalances != nil:
		balances := q.GetBalance(request.AllBalances.Address)
		return &AllBalancesResponse{Amount: balances}, nil
	default:
		return nil, NewError("Unknown BankQuery variant")
	}
}

func (q *MockQuerier) SetBalance(addr string, balance []Coin) {
	// clone coins so we don't accidentally edit them
	var empty []Coin
	q.Balances[addr] = append(empty, balance...)
}

func (q *MockQuerier) GetBalance(addr string) []Coin {
	bal := q.Balances[addr]
	if len(bal) == 0 {
		return bal
	}
	// if not empty, clone data
	var empty []Coin
	return append(empty, bal...)
}

func DisplayMessage(msg []byte) int {
	fmt.Println("Logging" + string(msg))
	return 0
}
