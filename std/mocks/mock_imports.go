package mocks

import (
	"errors"
	"fmt"

	"github.com/cosmwasm/cosmwasm-go/std/math"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

const (
	// ContractAddress is the default contract address returned by Env.
	ContractAddress = "test-contract"
	// BlockHeight is the default height returned by Env.
	BlockHeight = 12_345
	// BlockTime is the default block time returned by Env.
	BlockTime = 1_571_797_419_404_808_777
	// ChainID is the default chain ID returned by Env.
	ChainID = "cosmos-testnet-14002"
)

const (
	canonicalAddressLength = 32
)

var (
	_ std.Iterator        = (*Iterator)(nil)
	_ std.ReadonlyStorage = (*Storage)(nil)
	_ std.Storage         = (*Storage)(nil)
	_ std.Querier         = (*Querier)(nil)
	_ std.Api             = (*API)(nil)
)

// Deps returns mocked dependencies, funds can be provided optionally.
func Deps(funds []types.Coin) *std.Deps {
	return &std.Deps{
		Storage: NewStorage(),
		Api:     API{},
		Querier: NewQuerier(funds),
	}
}

// Env returns mocked environment.
func Env() types.Env {
	return types.Env{
		Block: types.BlockInfo{
			Height:  BlockHeight,
			Time:    BlockTime,
			ChainID: ChainID,
		},
		Contract: types.ContractInfo{
			Address: ContractAddress,
		},
	}
}

// Info returns mocked message info, given a sender and the funds.
func Info(sender string, funds []types.Coin) types.MessageInfo {
	return types.MessageInfo{
		Sender: sender,
		Funds:  funds,
	}
}

// Iterator mocks the std.Iterator.
type Iterator struct {
	Iter dbm.Iterator
}

func newIterator(iter dbm.Iterator) Iterator {
	return Iterator{
		Iter: iter,
	}
}

func (i Iterator) Next() (key, value []byte, err error) {
	if !i.Iter.Valid() {
		i.Iter.Close()
		return key, value, std.ErrIteratorDone
	}
	key, value = i.Iter.Key(), i.Iter.Value()
	i.Iter.Next()
	return
}

type Storage struct {
	storage dbm.DB
}

func NewStorage() *Storage {
	return &Storage{
		storage: dbm.NewMemDB(),
	}
}

func (s *Storage) Get(key []byte) []byte {
	v, err := s.storage.Get(key)
	if err != nil {
		// tm-db says that if the key is not found then the
		// value is nil, so we can panic here.
		panic(err)
	}

	return v
}

func (s *Storage) Range(start, end []byte, order std.Order) (iter std.Iterator) {
	var (
		iterator dbm.Iterator
		err      error
	)

	switch order {
	case std.Ascending:
		iterator, err = s.storage.Iterator(start, end)
		iter = newIterator(iterator)
	case std.Descending:
		iterator, err = s.storage.ReverseIterator(start, end)
		iter = newIterator(iterator)
	default:
		err = errors.New("unexpected Order")
	}

	if err != nil {
		panic(err)
	}
	return
}

func (s *Storage) Set(key, value []byte) {
	err := s.storage.Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (s *Storage) Remove(key []byte) {
	err := s.storage.Delete(key)
	if err != nil {
		panic(err)
	}
}

type API struct{}

func (a API) CanonicalAddress(human string) (types.CanonicalAddress, error) {
	if len(human) == 0 {
		return nil, errors.New("empty address")
	}
	if len(human) > canonicalAddressLength {
		return nil, errors.New("human encoding too long")
	}

	return []byte(human), nil
}

func (a API) HumanAddress(canonical types.CanonicalAddress) (string, error) {
	if len(canonical) != canonicalAddressLength {
		return "", errors.New("wrong canonical address length")
	}

	cutIndex := canonicalAddressLength
	for i, v := range canonical {
		if v == 0 {
			cutIndex = i
			break
		}
	}

	return string(canonical[:cutIndex]), nil
}

func (a API) ValidateAddress(human string) error {
	if len(human) > canonicalAddressLength {
		return errors.New("human encoding too long")
	}
	return nil
}

func (a API) Debug(msg string) {
	fmt.Println("DEBUG: " + msg)
}

type Querier struct {
	Balances map[string][]types.Coin
}

func NewQuerier(funds []types.Coin) *Querier {
	q := Querier{
		Balances: make(map[string][]types.Coin),
	}
	if len(funds) > 0 {
		q.SetBalance(ContractAddress, funds)
	}
	return &q
}

func (q *Querier) RawQuery(raw []byte) ([]byte, error) {
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

func (q *Querier) HandleQuery(request types.QueryRequest) (std.JSONType, error) {
	switch {
	case request.Bank != nil:
		return q.HandleBank(request.Bank)
	case request.Staking != nil:
		return nil, errors.New("staking queries not implemented")
	case request.Wasm != nil:
		return nil, errors.New("wasm queries not implemented")
	case request.Custom != nil:
		return nil, errors.New("custom queries not implemented")
	default:
		return nil, errors.New("unknown types.QueryRequest variant")
	}
}

func (q *Querier) HandleBank(request *types.BankQuery) (std.JSONType, error) {
	switch {
	case request.Balance != nil:
		balances := q.GetBalance(request.Balance.Address)
		coin := types.Coin{Denom: request.Balance.Denom, Amount: math.ZeroUint128()}
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
		return nil, errors.New("unknown types.BankQuery variant")
	}
}

func (q *Querier) SetBalance(addr string, balance []types.Coin) {
	// clone coins so we don't accidentally edit them
	var empty []types.Coin
	q.Balances[addr] = append(empty, balance...)
}

func (q *Querier) GetBalance(addr string) []types.Coin {
	bal := q.Balances[addr]
	if len(bal) == 0 {
		return bal
	}
	// if not empty, clone data
	var empty []types.Coin
	return append(empty, bal...)
}
