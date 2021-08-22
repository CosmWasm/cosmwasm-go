// +build !cosmwasm

package mocks

import (
	"fmt"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

func MockDeps(funds []types.Coin) *std.Deps {
	return &std.Deps{
		Storage: NewMockStorage(),
		Api:     MockApi{},
		Querier: NewMockQuerier(funds),
	}
}

const MOCK_CONTRACT_ADDR = "test-contract"

func MockEnv() types.Env {
	return types.Env{
		Block: types.BlockInfo{
			Height:  12_345,
			Time:    1_571_797_419_404_808_777,
			ChainID: "cosmos-testnet-14002",
		},
		Contract: types.ContractInfo{
			Address: MOCK_CONTRACT_ADDR,
		},
	}
}

func MockInfo(sender string, funds []types.Coin) types.MessageInfo {
	return types.MessageInfo{
		Sender: sender,
		Funds:  funds,
	}
}

var (
	_ std.Iterator = (*MockIterator)(nil)
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
		return key, value, types.NewError("the end of iterator")
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
	_ std.ReadonlyStorage = (*MockStorage)(nil)
	_ std.Storage         = (*MockStorage)(nil)
)

func (s *MockStorage) Get(key []byte) ([]byte, error) {
	return s.storage.Get(key)
}

func (s *MockStorage) Range(start, end []byte, order std.Order) (iter std.Iterator, err error) {
	var iterator dbm.Iterator
	switch order {
	case std.Ascending:
		iterator, err = s.storage.Iterator(start, end)
		iter = newMockIterator(iterator)
	case std.Descending:
		iterator, err = s.storage.ReverseIterator(start, end)
		iter = newMockIterator(iterator)
	default:
		err = types.NewError("failed. unexpected Order")
	}
	return
}

func (s *MockStorage) Set(key, value []byte) error {
	return s.storage.Set(key, value)
}

func (s *MockStorage) Remove(key []byte) error {
	return s.storage.Delete(key)
}

const canonicalLength = 32

// ensure Api interface compliance at compile time
var (
	_ std.Api = (*MockApi)(nil)
)

type MockApi struct{}

func (api MockApi) CanonicalAddress(human string) (types.CanonicalAddr, error) {
	if len(human) > canonicalLength {
		return nil, types.NewError("failed. human encoding too long")
	}

	return []byte(human), nil
}

func (api MockApi) HumanAddress(canonical types.CanonicalAddr) (string, error) {
	if len(canonical) != canonicalLength {
		return "", types.NewError("failed. wrong canonical address length")
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

func (api MockApi) ValidateAddress(human string) error {
	if len(human) > canonicalLength {
		return types.NewError("failed. human encoding too long")
	}
	return nil
}

func (api MockApi) Debug(msg string) {
	fmt.Println("DEBUG: " + msg)
}

// ====== Querier ======

// ensure Api interface compliance at compile time
var (
	_ std.Querier = (*MockQuerier)(nil)
)

type MockQuerier struct {
	Balances map[string][]types.Coin
}

func NewMockQuerier(funds []types.Coin) *MockQuerier {
	q := MockQuerier{
		Balances: make(map[string][]types.Coin),
	}
	if len(funds) > 0 {
		q.SetBalance(MOCK_CONTRACT_ADDR, funds)
	}
	return &q
}

func (q *MockQuerier) RawQuery(raw []byte) ([]byte, error) {
	var request types.QueryRequest
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

func (q *MockQuerier) HandleQuery(request types.QueryRequest) (std.JSONType, error) {
	switch {
	case request.Bank != nil:
		return q.HandleBank(request.Bank)
	case request.Staking != nil:
		return nil, types.NewError("Staking queries not implemented")
	case request.Wasm != nil:
		return nil, types.NewError("Wasm queries not implemented")
	case request.Custom != nil:
		return nil, types.NewError("Custom queries not implemented")
	default:
		return nil, types.NewError("Unknown types.QueryRequest variant")
	}
}

func (q *MockQuerier) HandleBank(request *types.BankQuery) (std.JSONType, error) {
	switch {
	case request.Balance != nil:
		balances := q.GetBalance(request.Balance.Address)
		coin := types.Coin{Denom: request.Balance.Denom, Amount: "0"}
		for _, c := range balances {
			if c.Denom == coin.Denom {
				coin.Amount = c.Amount
				break
			}
		}
		return &types.BalanceResponse{Amount: coin}, nil
	case request.AllBalances != nil:
		balances := q.GetBalance(request.AllBalances.Address)
		return &types.AllBalancesResponse{Amount: balances}, nil
	default:
		return nil, types.NewError("Unknown types.BankQuery variant")
	}
}

func (q *MockQuerier) SetBalance(addr string, balance []types.Coin) {
	// clone coins so we don't accidentally edit them
	var empty []types.Coin
	q.Balances[addr] = append(empty, balance...)
}

func (q *MockQuerier) GetBalance(addr string) []types.Coin {
	bal := q.Balances[addr]
	if len(bal) == 0 {
		return bal
	}
	// if not empty, clone data
	var empty []types.Coin
	return append(empty, bal...)
}

func DisplayMessage(msg []byte) int {
	fmt.Println("Logging" + string(msg))
	return 0
}
