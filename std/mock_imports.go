// +build !cosmwasm

package std

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
	IteratorId uint32
}

func (iterator ExternalIterator) Next() (key, value []byte, err error) {

	return key, value, nil
}

var storage map[string][]byte

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

func (storage ExternalStorage) Get(key []byte) (value []byte, err error) {

	return nil, nil
}

func (storage ExternalStorage) Range(start, end []byte, order Order) (Iterator, error) {
	return nil, nil
}

func (storage ExternalStorage) Set(key, value []byte) error {

	return nil
}

func (storage ExternalStorage) Remove(key []byte) error {

	return nil
}

type CanonicalAddr []byte

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

	return CanonicalAddr{}, nil
}

func (api ExternalApi) HumanAddress(canonical CanonicalAddr) (string, error) {

	return string(""), nil
}

// ====== Querier ======
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
