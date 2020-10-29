package std

// =========== Extern --> context =======
type Extern struct {
	EStorage Storage
	EApi     Api
	EQuerier Querier
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
	Debug(msg string)
}

type Querier interface {
	RawQuery(request []byte) ([]byte, error)
}
