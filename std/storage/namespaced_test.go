package storage

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO(fdymylja): test this

func TestNamespaced(t *testing.T) {
	storage := mock.Storage()
	ns := NewNamespaced("prefix")

	ns.Set(storage, []byte("1"), []byte("a"))
	ns.Set(storage, []byte("11"), []byte("aa"))

	iter := ns.Range(storage, []byte("1"), nil, std.Ascending)
	// test iterator strips prefix bytes
	k, v, err := iter.Next()
	require.NoError(t, err)
	require.Equal(t, []byte("1"), k)
	require.Equal(t, []byte("a"), v)

	k, v, err = iter.Next()
	require.NoError(t, err)
	require.Equal(t, []byte("11"), k)
	require.Equal(t, []byte("aa"), v)

	// test get
	v = ns.Get(storage, []byte("1"))
	require.Equal(t, []byte("a"), v)

	// test remove
	ns.Remove(storage, []byte("1"))
	require.Nil(t, storage.Get([]byte("1")))

}
