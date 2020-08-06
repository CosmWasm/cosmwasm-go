package std

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExternalStorage(t *testing.T) {
	var es ExternalStorage
	key1, key2, key3, key4, key5 := []byte("aaaaa"), []byte("bbbbb"), []byte("ccccc"), []byte("ddddd"), []byte("eeeee")
	value1, value2, value3, value4, value5 := []byte("11111"), []byte("22222"), []byte("33333"), []byte("44444"), []byte("55555")

	// setter && getter
	bytes, err := es.Get(key1)
	require.Nil(t, bytes)
	require.NoError(t, err)
	require.NoError(t, es.Set(key1, value1))
	require.NoError(t, es.Set(key2, value2))
	require.NoError(t, es.Set(key3, value3))
	require.NoError(t, es.Set(key4, value4))
	require.NoError(t, es.Set(key5, value5))

	// iterator
	// ascending
	iter, err := es.Range([]byte{'a'}, []byte{'d'}, Ascending)
	require.NoError(t, err)
	assertKV(t, iter, key1, value1, false)
	assertKV(t, iter, key2, value2, false)
	assertKV(t, iter, key3, value3, false)
	assertKV(t, iter, key4, value4, true)
	// descending
	iter, err = es.Range([]byte{'b'}, []byte("eeeef"), Descending)
	assertKV(t, iter, key5, value5, false)
	assertKV(t, iter, key4, value4, false)
	assertKV(t, iter, key3, value3, false)
	assertKV(t, iter, key2, value2, false)
	assertKV(t, iter, key1, value1, true)
}

func assertKV(t *testing.T, iter Iterator, key, value []byte, isEnd bool) {
	curKey, curValue, err := iter.Next()
	if isEnd {
		require.Error(t, err)
		require.Nil(t, curKey)
		require.Nil(t, curValue)
		return
	}
	require.NoError(t, err)
	require.Equal(t, curKey, key)
	require.Equal(t, curValue, value)
}
