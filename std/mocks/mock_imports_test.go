//go:build !cosmwasm
// +build !cosmwasm

package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

func TestMockStorage(t *testing.T) {
	es := NewMockStorage()
	key1, key2, key3, key4, key5 := []byte("aaaaa"), []byte("bbbbb"), []byte("ccccc"), []byte("ddddd"), []byte("eeeee")
	value1, value2, value3, value4, value5 := []byte("11111"), []byte("22222"), []byte("33333"), []byte("44444"), []byte("55555")
	inexistentKey := []byte("inexistent")

	// setter && getter
	bytes := es.Get(key1)
	require.Nil(t, bytes)
	es.Set(key1, value1)
	es.Set(key2, value2)
	es.Set(key3, value3)
	es.Set(key4, value4)
	es.Set(key5, value5)

	// iterator
	// ascending
	iter := es.Range([]byte{'a'}, []byte{'d'}, std.Ascending)
	assertKV(t, iter, key1, value1, false)
	assertKV(t, iter, key2, value2, false)
	assertKV(t, iter, key3, value3, false)
	assertKV(t, iter, key4, value4, true)
	// descending
	iter = es.Range([]byte{'b'}, []byte("eeeef"), std.Descending)
	assertKV(t, iter, key5, value5, false)
	assertKV(t, iter, key4, value4, false)
	assertKV(t, iter, key3, value3, false)
	assertKV(t, iter, key2, value2, false)
	assertKV(t, iter, key1, value1, true)

	// delete
	es.Remove(inexistentKey)
	es.Remove(key1)
	bytes = es.Get(key1)
	require.Nil(t, bytes)

}

func assertKV(t *testing.T, iter std.Iterator, key, value []byte, isEnd bool) {
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

func TestMockApi_CanonicalAddress(t *testing.T) {
	ea := MockApi{}
	humanAddr := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	longHumanAddr := humanAddr + "a"
	expectedCanonAddr := types.CanonicalAddress(humanAddr)

	canonAddr, err := ea.CanonicalAddress(longHumanAddr)
	require.Error(t, err)
	require.Nil(t, canonAddr)

	canonAddr, err = ea.CanonicalAddress(humanAddr)
	require.NoError(t, err)
	require.Equal(t, expectedCanonAddr, canonAddr)
}

func TestMockApi_HumanAddress(t *testing.T) {
	ea := MockApi{}
	expectedHumanAddr := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	expectedCanonAddr := types.CanonicalAddress(expectedHumanAddr)

	humanAddr, err := ea.HumanAddress(expectedCanonAddr)
	require.NoError(t, err)
	require.Equal(t, expectedHumanAddr, humanAddr)

	// error report
	longCanonAddr := make(types.CanonicalAddress, canonicalLength)
	copy(longCanonAddr, expectedCanonAddr)
	longCanonAddr = append(longCanonAddr, 'a')
	humanAddr, err = ea.HumanAddress(longCanonAddr)
	require.Error(t, err)
	require.Equal(t, "", humanAddr)

	inputCanonAddr := make(types.CanonicalAddress, canonicalLength)
	copy(inputCanonAddr, expectedCanonAddr)
	inputCanonAddr[9] = 0
	humanAddr, err = ea.HumanAddress(inputCanonAddr)
	require.NoError(t, err)
	require.Equal(t, "aaaaaaaaa", humanAddr)
}
