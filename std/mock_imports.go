// +build !cosmwasm

package std

import (
	"errors"
	"fmt"
	dbm "github.com/tendermint/tm-db"
)

func MockExtern() *Extern {
	return &Extern{
		EStorage: NewMockStorage(),
		EApi:     MockApi{},
		EQuerier: MockQuerier{},
	}
}

const MOCK_CONTRACT_ADDR = "test-contract"

func MockEnv(sender string, funds []Coin) Env {
	return Env{
		Block: BlockInfo{
			Height:  12_345,
			Time:    1_571_797_419,
			ChainID: "cosmos-testnet-14002",
		},
		Message: MessageInfo{
			Sender:    sender,
			SentFunds: funds,
		},
		Contract: ContractInfo{
			Address: MOCK_CONTRACT_ADDR,
		},
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
		return key, value, errors.New("the end of iterator")
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
		err = errors.New("failed. unexpected Order")
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
		return nil, errors.New("failed. human encoding too long")
	}

	return []byte(human), nil
}

func (api MockApi) HumanAddress(canonical CanonicalAddr) (string, error) {
	if len(canonical) != canonicalLength {
		return "", errors.New("failed. wrong canonical address length")
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

// ====== Querier ======

// ensure Api interface compliance at compile time
var (
	_ Querier = (*MockQuerier)(nil)
)

type MockQuerier struct{}

func (querier MockQuerier) RawQuery(request []byte) ([]byte, error) {
	return []byte(""), nil
}

func DisplayMessage(msg []byte) int {
	fmt.Println("Logging" + string(msg))
	return 0
}
