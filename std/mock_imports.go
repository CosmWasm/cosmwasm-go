// +build !cosmwasm

package std

import (
	"errors"
	"fmt"
	dbm "github.com/tendermint/tm-db"
)

// ====== DB mock ======
type Order uint32

const (
	Ascending  Order = 1
	Descending Order = 2
)

type Iterator interface {
	Next() (key, value []byte, err error)
}

var (
	_ Iterator = (*ExternalIterator)(nil)
)

type ExternalIterator struct {
	Iter dbm.Iterator
}

func newExternalIterator(iter dbm.Iterator) ExternalIterator {
	return ExternalIterator{
		Iter: iter,
	}
}

func (iter ExternalIterator) Next() (key, value []byte, err error) {
	if !iter.Iter.Valid() {
		iter.Iter.Close()
		return key, value, errors.New("the end of iterator")
	}
	key, value = iter.Iter.Key(), iter.Iter.Value()
	iter.Iter.Next()
	return
}

var storage = dbm.NewMemDB()

type ReadonlyStorage interface {
	Get(key []byte) (value []byte, err error)
	Range(start, end []byte, order Order) (Iterator, error)
}

type Storage interface {
	ReadonlyStorage

	Set(key, value []byte) error
	Remove(key []byte) error
}

type ExternalStorage struct{}

var (
	_ ReadonlyStorage = (*ExternalStorage)(nil)
	_ Storage         = (*ExternalStorage)(nil)
)

func (es ExternalStorage) Get(key []byte) ([]byte, error) {
	return storage.Get(key)
}

func (es ExternalStorage) Range(start, end []byte, order Order) (iter Iterator, err error) {
	var iterator dbm.Iterator
	switch order {
	case Ascending:
		iterator, err = storage.Iterator(start, end)
		iter = newExternalIterator(iterator)
	case Descending:
		iterator, err = storage.ReverseIterator(start, end)
		iter = newExternalIterator(iterator)
	default:
		err = errors.New("failed. unexpected Order")
	}
	return
}

func (es ExternalStorage) Set(key, value []byte) error {
	return storage.Set(key, value)
}

func (es ExternalStorage) Remove(key []byte) error {
	return storage.Delete(key)
}

type CanonicalAddr []byte

const canonicalLength = 32

type Api interface {
	CanonicalAddress(human string) (CanonicalAddr, error)
	HumanAddress(canonical CanonicalAddr) (string, error)
}

// ensure Api interface compliance at compile time
var (
	_ Api = (*ExternalApi)(nil)
)

type ExternalApi struct{}

func (api ExternalApi) CanonicalAddress(human string) (CanonicalAddr, error) {
	if len(human) > canonicalLength {
		return nil, errors.New("failed. human encoding too long")
	}

	return []byte(human), nil
}

func (api ExternalApi) HumanAddress(canonical CanonicalAddr) (string, error) {
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

// ------- query detail types ---------
type QueryResponseOk struct {
	Ok []byte `json:"Ok,omitempty"`
}

// This is a 2-level result
type QuerierResult struct {
	Ok QueryResponseOk `json:"Ok,omitempty"`
}

type Querier interface {
	RawQuery(request []byte) ([]byte, error)
}

// ensure Api interface compliance at compile time
var (
	_ Querier = (*ExternalQuerier)(nil)
)

type ExternalQuerier struct{}

func (querier ExternalQuerier) RawQuery(request []byte) ([]byte, error) {

	return []byte(""), nil
}

func DisplayMessage(msg []byte) int {
	fmt.Println("Logging" + string(msg))
	return 0
}
