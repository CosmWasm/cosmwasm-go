// +build cosmwasm

package std

/*
#include "stdlib.h"
extern void* db_read(void* key);
extern void db_write(void* key, void* value);
extern void db_remove(void* key);

extern int db_scan(void* start_ptr, void* end_ptr, int order);
extern void* db_next(unsigned iterator_id);

extern int canonicalize_address(void* human, void* canonical);
extern int humanize_address(void* canonical, void* human);
extern void debug(void* msg);

extern void* query_chain(void* request);

extern int display_message(void* str);
*/
import "C"

import (
	"unsafe"
)

const (
	// A kibi (kilo binary)
	KI uint32 = 1024

	// The number of bytes of the memory region we pre-allocate for the result data in ExternalIterator.next
	DB_READ_KEY_BUFFER_LENGTH uint32 = 64 * KI

	// The number of bytes of the memory region we pre-allocate for the result data in ExternalStorage.get and ExternalIterator.next
	DB_READ_VALUE_BUFFER_LENGTH uint32 = 128 * KI

	// The number of bytes of the memory region we pre-allocate for the result data in queries
	QUERY_RESULT_BUFFER_LENGTH uint32 = 128 * KI

	// An upper bound for typical canonical address lengths (e.g. 20 in Cosmos SDK/Ethereum or 32 in Nano/Substrate)
	CANONICAL_ADDRESS_BUFFER_LENGTH uint32 = 32

	// An upper bound for typical human readable address formats (e.g. 42 for Ethereum hex addresses or 90 for bech32)
	HUMAN_ADDRESS_BUFFER_LENGTH uint32 = 90
)

// log print

func DisplayMessage(data []byte) int {
	//cause that cosmwasm-vm does not support display_message, we denied it in publish version for now.
	//if you want using display_message to build and test your contract, try using cosmwasm-simulate tool to load and test
	//download it from :https://github.com/CosmWasm/cosmwasm-simulate

	msg := C.malloc(C.ulong(len(data)))
	regionMsg := TranslateToRegion(data, uintptr(msg))
	C.display_message(unsafe.Pointer(regionMsg))
	C.free(unsafe.Pointer(msg))
	return 0
}

// ====== DB ======

var (
	_ ReadonlyStorage = (*ExternalStorage)(nil)
	_ Storage         = (*ExternalStorage)(nil)
)

type ExternalStorage struct{}

func (storage ExternalStorage) Get(key []byte) (value []byte, err error) {
	keyPtr := C.malloc(C.ulong(len(key)))
	regionKey := TranslateToRegion(key, uintptr(keyPtr))

	read := C.db_read(unsafe.Pointer(regionKey))
	C.free(unsafe.Pointer(keyPtr))

	if read == nil {
		return nil, NewError("key not existed")
	}

	b := TranslateToSlice(uintptr(read))
	//maybe have memory leak
	return b, nil
}

func (storage ExternalStorage) Range(start, end []byte, order Order) (Iterator, error) {
	/*
		ptrStart := C.malloc(C.ulong(len(start)))
		regionStart := TranslateToRegion(start, uintptr(ptrStart))

		ptrEnd := C.malloc(C.ulong(len(end)))
		regionEnd := TranslateToRegion(end, uintptr(ptrEnd))

		iterId := C.db_scan(unsafe.Pointer(regionStart), unsafe.Pointer(regionEnd), C.int(order))
		C.free(ptrStart)
		C.free(ptrEnd)

		if iterId < 0 {
			return nil, NewError("error creating iterator (via db_scan): " + string(int(iterId)))
		}

		return ExternalIterator{uint32(iterId)}, nil
	*/
	return nil, nil
}

func (storage ExternalStorage) Set(key, value []byte) error {
	ptrKey := C.malloc(C.ulong(len(key)))
	ptrVal := C.malloc(C.ulong(len(value)))
	regionKey := TranslateToRegion(key, uintptr(ptrKey))
	regionValue := TranslateToRegion(value, uintptr(ptrVal))

	C.db_write(unsafe.Pointer(regionKey), unsafe.Pointer(regionValue))
	C.free(ptrKey)
	C.free(ptrVal)

	return nil
}

func (storage ExternalStorage) Remove(key []byte) error {
	keyPtr := C.malloc(C.ulong(len(key)))
	regionKey := TranslateToRegion(key, uintptr(keyPtr))

	C.db_remove(unsafe.Pointer(regionKey))
	C.free(keyPtr)

	return nil
}

var (
	_ Iterator = (*ExternalIterator)(nil)
)

type ExternalIterator struct {
	IteratorId uint32
}

func (iterator ExternalIterator) Next() (key, value []byte, err error) {
	/*
		regionKey, _ := Build_region(DB_READ_KEY_BUFFER_LENGTH, 0)
		regionNextValue, _ := Build_region(DB_READ_VALUE_BUFFER_LENGTH, 0)

		ret := nil //C.db_next(C.uint(iterator.IteratorId))

		if ret == nil {
			return nil, nil, NewError("unknown error from db_next ")
		}

		key = TranslateToSlice(uintptr(regionKey))
		value = TranslateToSlice(uintptr(regionNextValue))

		if len(key) == 0 {
			return nil, nil, NewError("empty key get from db_next")
		}

		return key, value, nil

	*/
	return nil, nil, NewError("unsupported for now")
}

// ====== API ======
type CanonicalAddr []byte

// ensure Api interface compliance at compile time
var (
	_ Api = (*ExternalApi)(nil)
)

type ExternalApi struct{}

func (api ExternalApi) CanonicalAddress(human string) (CanonicalAddr, error) {
	humanAddr := []byte(human)
	humanPtr := C.malloc(C.ulong(len(humanAddr)))
	regionHuman := TranslateToRegion(humanAddr, uintptr(humanPtr))

	regionCanon, _ := Build_region(CANONICAL_ADDRESS_BUFFER_LENGTH, 0)

	ret := C.canonicalize_address(unsafe.Pointer(regionHuman), unsafe.Pointer(regionCanon))
	C.free(humanPtr)

	if ret < 0 {
		return nil, NewError("canonicalize_address returned error")
	}

	canoAddress := TranslateToSlice(uintptr(regionCanon))

	return canoAddress, nil
}

func (api ExternalApi) HumanAddress(canonical CanonicalAddr) (string, error) {
	canonPtr := C.malloc(C.ulong(len(canonical)))
	regionCanon := TranslateToRegion(canonical, uintptr(canonPtr))

	regionHuman, _ := Build_region(HUMAN_ADDRESS_BUFFER_LENGTH, 0)

	ret := C.humanize_address(unsafe.Pointer(regionCanon), unsafe.Pointer(regionHuman))
	C.free(canonPtr)

	if ret < 0 {
		return "", NewError("humanize_address returned error")
	}

	humanAddress := TranslateToSlice(uintptr(regionHuman))

	return string(humanAddress), nil
}

func (api ExternalApi) Debug(msg string) {
	msgPtr := C.malloc(C.ulong(len(msg)))
	regionMsg := TranslateToRegion([]byte(msg), uintptr(msgPtr))
	C.debug(unsafe.Pointer(regionMsg))
	C.free(msgPtr)
}

// ====== Querier ======

// ensure Api interface compliance at compile time
var (
	_ Querier = (*ExternalQuerier)(nil)
)

type ExternalQuerier struct{}

func (querier ExternalQuerier) RawQuery(request []byte) ([]byte, error) {
	reqPtr := C.malloc(C.ulong(len(request)))
	regionReq := TranslateToRegion(request, uintptr(reqPtr))

	ret := C.query_chain(unsafe.Pointer(regionReq))
	C.free(reqPtr)

	if ret == nil {
		return nil, NewError("failed to query chain: unknown error")
	}

	response := TranslateToSlice(uintptr(ret))
	// TODO: parse this into the proper structure
	// success looks like: {"ok":{"ok":"eyJhbW91bnQiOlt7ImRlbm9tIjoid2VpIiwiYW1vdW50IjoiNzY1NDMyIn1dfQ=="}}
	var qres QuerierResult
	err := qres.UnmarshalJSON(response)
	if err != nil {
		return nil, err
	}
	if qres.Error != nil {
		return nil, NewError(qres.Error.Error())
	}
	if qres.Ok.Error != "" {
		return nil, NewError(qres.Ok.Error)
	}
	return qres.Ok.Ok, nil
}

// use for ezjson Logging
// TODO: remove????

func Wasmlog(msg []byte) int {
	msgPtr := C.malloc(C.ulong(len(msg)))
	regionMsg := TranslateToRegion(msg, uintptr(msgPtr))
	C.debug(unsafe.Pointer(regionMsg))
	C.free(msgPtr)
	return 0
}
