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

extern void* query_chain(void* request);

extern int display_message(void* str);
*/
import "C"

import (
	"encoding/base64"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
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
		return nil, errors.New("key not existed")
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
			return nil, errors.New("error creating iterator (via db_scan): " + string(int(iterId)))
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
	/*
		regionKey, _ := Build_region(DB_READ_KEY_BUFFER_LENGTH, 0)
		regionNextValue, _ := Build_region(DB_READ_VALUE_BUFFER_LENGTH, 0)

		ret := nil //C.db_next(C.uint(iterator.IteratorId))

		if ret == nil {
			return nil, nil, errors.New("unknown error from db_next ")
		}

		key = TranslateToSlice(uintptr(regionKey))
		value = TranslateToSlice(uintptr(regionNextValue))

		if len(key) == 0 {
			return nil, nil, errors.New("empty key get from db_next")
		}

		return key, value, nil

	*/
	return nil, nil, errors.New("unsupported for now")
}

// ====== API ======
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
	humanAddr := []byte(human)
	humanPtr := C.malloc(C.ulong(len(humanAddr)))
	regionHuman := TranslateToRegion(humanAddr, uintptr(humanPtr))

	regionCanon, _ := Build_region(CANONICAL_ADDRESS_BUFFER_LENGTH, 0)

	ret := C.canonicalize_address(unsafe.Pointer(regionHuman), unsafe.Pointer(regionCanon))
	C.free(humanPtr)

	if ret < 0 {
		return nil, errors.New("canonicalize_address returned error")
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
		return "", errors.New("humanize_address returned error")
	}

	humanAddress := TranslateToSlice(uintptr(regionHuman))

	return string(humanAddress), nil
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
	reqPtr := C.malloc(C.ulong(len(request)))
	regionReq := TranslateToRegion(request, uintptr(reqPtr))

	ret := C.query_chain(unsafe.Pointer(regionReq))
	C.free(reqPtr)

	if ret == nil {
		return nil, errors.New("failed to query chain: unknown error")
	}

	response := TranslateToSlice(uintptr(ret))
	return response, nil
}

// ------- query detail types ---------
type QueryResponseOk struct {
	Ok string `json:"Ok,omitempty,rust_option"`
}

// This is a 2-level result
type QuerierResult struct {
	Ok QueryResponseOk `json:"Ok,omitempty"`
}

func BuildQueryResponse(msg string) *QueryResponseOk {
	encoding := base64.StdEncoding.EncodeToString([]byte(msg))
	return &QueryResponseOk{Ok: encoding}
}

// QueryRequest is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type QueryRequest struct {
	Bank    *BankQuery    `json:"bank,omitempty"`
	Custom  RawMessage    `json:"custom,omitempty"`
	Staking *StakingQuery `json:"staking,omitempty"`
	Wasm    *WasmQuery    `json:"wasm,omitempty"`
}

type BankQuery struct {
	Balance     *BalanceQuery     `json:"balance,omitempty"`
	AllBalances *AllBalancesQuery `json:"all_balances,omitempty"`
}

type BalanceQuery struct {
	Address string `json:"address"`
	Denom   string `json:"denom"`
}

// BalanceResponse is the expected response to BalanceQuery
type BalanceResponse struct {
	Amount Coin `json:"amount"`
}

type AllBalancesQuery struct {
	Address string `json:"address"`
}

// AllBalancesResponse is the expected response to AllBalancesQuery
type AllBalancesResponse struct {
	Amount Coins `json:"amount"`
}

type StakingQuery struct {
	Validators     ValidatorsQuery     `json:"validators,omitempty"`
	AllDelegations AllDelegationsQuery `json:"all_delegations,omitempty"`
	Delegation     DelegationQuery     `json:"delegation,omitempty"`
	BondedDenom    struct{}            `json:"bonded_denom,omitempty"`
}

type ValidatorsQuery struct{}

// ValidatorsResponse is the expected response to ValidatorsQuery
type ValidatorsResponse struct {
	Validators Validators `json:"validators"`
}

// TODO: Validators must JSON encode empty array as []
type Validators []Validator

// MarshalJSON ensures that we get [] for empty arrays
func (v Validators) MarshalJSON() ([]byte, error) {
	if len(v) == 0 {
		return []byte("[]"), nil
	}
	var raw []Validator = v
	return ezjson.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (v *Validators) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Validator
	if err := ezjson.Unmarshal(data, &raw); err != nil {
		return err
	}
	*v = raw
	return nil
}

type Validator struct {
	Address string `json:"address"`
	// decimal string, eg "0.02"
	Commission string `json:"commission"`
	// decimal string, eg "0.02"
	MaxCommission string `json:"max_commission"`
	// decimal string, eg "0.02"
	MaxChangeRate string `json:"max_change_rate"`
}

type AllDelegationsQuery struct {
	Delegator string `json:"delegator"`
}

type DelegationQuery struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
}

// AllDelegationsResponse is the expected response to AllDelegationsQuery
type AllDelegationsResponse struct {
	Delegations Delegations `json:"delegations"`
}

type Delegations []Delegation

// MarshalJSON ensures that we get [] for empty arrays
func (d Delegations) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	var raw []Delegation = d
	return ezjson.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (d *Delegations) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Delegation
	if err := ezjson.Unmarshal(data, &raw); err != nil {
		return err
	}
	*d = raw
	return nil
}

type Delegation struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

// DelegationResponse is the expected response to DelegationsQuery
type DelegationResponse struct {
	Delegation *FullDelegation `json:"delegation,omitempty"`
}

type FullDelegation struct {
	Delegator          string `json:"delegator"`
	Validator          string `json:"validator"`
	Amount             Coin   `json:"amount"`
	AccumulatedRewards Coin   `json:"accumulated_rewards"`
	CanRedelegate      Coin   `json:"can_redelegate"`
}

type BondedDenomResponse struct {
	Denom string `json:"denom"`
}

type WasmQuery struct {
	Smart *SmartQuery `json:"smart,omitempty"`
	Raw   *RawQuery   `json:"raw,omitempty"`
}

// SmartQuery respone is raw bytes ([]byte)
type SmartQuery struct {
	ContractAddr string `json:"contract_addr"`
	Msg          []byte `json:"msg"`
}

// RawQuery response is raw bytes ([]byte)
type RawQuery struct {
	ContractAddr string `json:"contract_addr"`
	Key          []byte `json:"key"`
}
