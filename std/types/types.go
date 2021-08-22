package types

import "strconv"

type CanonicalAddr []byte

// Coin is a string representation of the sdk.Coin type (more portable than sdk.Int)
type Coin struct {
	Denom  string // type, eg. "ATOM"
	Amount string // string encoing of decimal value, eg. "12.3456"
}

func NewCoin(amount uint64, denom string) Coin {
	return Coin{
		Denom:  denom,
		Amount: strconv.FormatUint(amount, 10),
	}
}

func NewCoins(amount uint64, denom string) []Coin {
	return []Coin{NewCoin(amount, denom)}
}

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return GenericError("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
