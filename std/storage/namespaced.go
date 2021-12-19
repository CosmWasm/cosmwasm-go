package storage

import (
	"github.com/cosmwasm/cosmwasm-go/std"
)

const (
	// NamespaceTerminator is the uint8 that terminates a namespace.
	NamespaceTerminator = 0x00
)

// Namespaced defines a higher level storage layer that saves bytes in a prefix.
type Namespaced struct {
	prefix []byte
}

// NewNamespaced generates a new Namespaced instance given a unique namespace.
func NewNamespaced(namespace string) Namespaced {
	return Namespaced{
		prefix: append([]byte(namespace), NamespaceTerminator),
	}
}

// Set sets key and value after prepending the namespace prefix.
func (s Namespaced) Set(storage std.Storage, key, value []byte) {
	storage.Set(s.key(key), value)
}

// Get gets the given key after prepending the namespace prefix.
func (s Namespaced) Get(storage std.Storage, key []byte) []byte {
	return storage.Get(s.key(key))
}

// Remove removes the given key after prepending the namespace prefix.
func (s Namespaced) Remove(storage std.Storage, key []byte) {
	storage.Remove(s.key(key))
}

// Range ranges over the given start to end range after prepending the namespace prefix.
func (s Namespaced) Range(storage std.Storage, start, end []byte, order std.Order) std.Iterator {
	nsStart := s.key(start)
	nsEnd := s.rangeEndKey(end)
	iter := storage.Range(nsStart, nsEnd, order)

	return NamespacedIterator{prefixLength: len(s.prefix), iter: iter}
}

// key returns the key with the namespace prefix.
func (s Namespaced) key(k []byte) []byte {
	return append(s.prefix, k...)
}

// rangeEndKey does special handling for a nil ending key.
// if the end is nil, then we take the prefix and increase it by one.
func (s Namespaced) rangeEndKey(end []byte) []byte {
	if end != nil {
		return s.key(end)
	}

	e := make([]byte, len(s.prefix))
	copy(e, s.prefix)

	for i := len(s.prefix) - 1; i >= 0; i-- {
		if e[i] < byte(0xFF) {
			e[i]++
			return e
		}
		e[i] = byte(0x00)
		if i == 0 {
			return nil
		}
	}

	return nil
}

// NamespacedIterator is the std.Iterator of namespaced storage.
type NamespacedIterator struct {
	prefixLength int
	iter         std.Iterator
}

// Next returns the next key value combo in the iterator, it removes the namespace prefix from key.
func (n NamespacedIterator) Next() (key, value []byte, err error) {
	rawKey, value, err := n.iter.Next()
	if err != nil {
		return nil, nil, err
	}

	return rawKey[n.prefixLength:], value, nil
}