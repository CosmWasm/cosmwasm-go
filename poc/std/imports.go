package std

/*
#include "stdlib.h"
extern int db_read(void* key, void* value);
extern int db_write(void* key, void* value);
extern int db_remove(void* key);
extern int db_scan(void* start_ptr, void* end_ptr, void* order, void* iterator_id);
extern int db_next(void* iterator_id, void* next_result);
extern int canonicalize_address(void* source, void* destination, void* result);
extern int humanize_address(void* source, void* destination, void* result);
extern int query_chain(void* request, void* response);
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
	"unsafe"
)

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
	regionValue, _ := Build_region(1024, 0)

	ret := C.db_read(unsafe.Pointer(regionKey), unsafe.Pointer(regionValue))
	C.free(keyPtr)

	if ret != 0 {
		if ret == -1001001 {
			return nil, errors.New("key not existed")
		}
		return nil, errors.New("call success but reading failed")
	}

	b := TranslateToSlice(uintptr(regionValue))
	return b, nil
}

func (storage ExternalStorage) Range(start, end []byte, order Order) (Iterator, error) {
	orderByte := bytes.NewBuffer([]byte{})
	binary.Write(orderByte, binary.BigEndian, &order)

	ptrStart := C.malloc(C.ulong(len(start)))
	regionStart := TranslateToRegion(start, uintptr(ptrStart))

	ptrEnd := C.malloc(C.ulong(len(end)))
	regionEnd := TranslateToRegion(end, uintptr(ptrEnd))

	ptrOrder := C.malloc(C.ulong(len(orderByte.Bytes())))
	regionOrder := TranslateToRegion(orderByte.Bytes(), uintptr(ptrOrder))

	regionIteratorId, _ := Build_region(256, 0)

	ret := C.db_scan(unsafe.Pointer(regionStart), unsafe.Pointer(regionEnd), unsafe.Pointer(&regionOrder), unsafe.Pointer(regionIteratorId))
	C.free(ptrStart)
	C.free(ptrEnd)
	C.free(ptrOrder)

	if ret != 0 {
		return nil, errors.New("some Error happens during scan storage")
	}

	iterId := TranslateToUint32(uintptr(regionIteratorId))
	if iterId == 0 {
		return nil, errors.New("failed to convert iterator id")
	}

	return ExternalIterator{iterId}, nil
}

func (storage ExternalStorage) Set(key, value []byte) error {
	ptrKey := C.malloc(C.ulong(len(key)))
	ptrVal := C.malloc(C.ulong(len(value)))
	regionKey := TranslateToRegion(key, uintptr(ptrKey))
	regionValue := TranslateToRegion(value, uintptr(ptrVal))

	ret := C.db_write(unsafe.Pointer(regionKey), unsafe.Pointer(regionValue))
	C.free(ptrKey)
	C.free(ptrVal)

	if ret != 0 {
		return errors.New("some Error happens during Write storage")
	}

	return nil
}

func (storage ExternalStorage) Remove(key []byte) error {
	keyPtr := C.malloc(C.ulong(len(key)))
	regionKey := TranslateToRegion(key, uintptr(keyPtr))

	ret := C.db_remove(unsafe.Pointer(regionKey))
	C.free(keyPtr)

	if ret != 0 {
		return errors.New("some Error happens during modify storage")
	}

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
	iterIdByte := bytes.NewBuffer([]byte{})
	binary.Write(iterIdByte, binary.BigEndian, &iterator.IteratorId)

	ptrIterId := C.malloc(C.ulong(len(iterIdByte.Bytes())))
	regionIterId := TranslateToRegion(iterIdByte.Bytes(), uintptr(ptrIterId))
	regionNextResult, _ := Build_region(1024, 0)

	ret := C.db_next(unsafe.Pointer(regionIterId), unsafe.Pointer(regionNextResult))
	C.free(ptrIterId)

	if ret != 0 {
		return nil, nil, errors.New("failed to get db next value")
	}

	kv := TranslateToSlice(uintptr(regionNextResult))

	bytesBuffer := bytes.NewBuffer(kv)
	var keyLen uint32
	err1 := binary.Read(bytesBuffer, binary.BigEndian, &keyLen)
	if err1 != nil || keyLen == 0 {
		return nil, nil, errors.New("failed to convert key len from db_next result")
	}

	kv = kv[:len(kv)-4]

	return kv[:len(kv)-int(keyLen)], kv[len(kv)-int(keyLen)+1:], nil
}

// ====== API ======
const (
	CANONICAL_ADDRESS_BUFFER_LENGTH uint32 = 32
	HUMAN_ADDRESS_BUFFER_LENGTH     uint32 = 90
)

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
	regionResult, _ := Build_region(256, 0)

	ret := C.canonicalize_address(unsafe.Pointer(regionHuman), unsafe.Pointer(regionCanon), unsafe.Pointer(regionResult))
	C.free(humanPtr)

	if ret != 0 {
		return nil, errors.New("fail to convert human address to canonicalize address")
	}

	result := TranslateToUint32(uintptr(regionResult))
	if result != 0 {
		return nil, errors.New("fail to convert human address to canonicalize address")
	}

	canoAddress := TranslateToSlice(uintptr(regionCanon))

	return canoAddress, nil
}

func (api ExternalApi) HumanAddress(canonical CanonicalAddr) (string, error) {
	canonPtr := C.malloc(C.ulong(len(canonical)))
	regionCanon := TranslateToRegion(canonical, uintptr(canonPtr))

	regionHuman, _ := Build_region(HUMAN_ADDRESS_BUFFER_LENGTH, 0)
	regionResult, _ := Build_region(256, 0)

	ret := C.canonicalize_address(unsafe.Pointer(regionCanon), unsafe.Pointer(regionHuman), unsafe.Pointer(regionResult))
	C.free(canonPtr)

	if ret != 0 {
		return "", errors.New("fail to convert human address to canonicalize address")
	}

	result := TranslateToUint32(uintptr(regionResult))
	if result != 0 {
		return "", errors.New("fail to convert human address to canonicalize address")
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

const (
	QueryChainBufferSize uint32 = 20 * 1024 * 1024
)

type ExternalQuerier struct{}

func (querier ExternalQuerier) RawQuery(request []byte) ([]byte, error) {
	reqPtr := C.malloc(C.ulong(len(request)))
	regionReq := TranslateToRegion(request, uintptr(reqPtr))
	regionResponse, _ := Build_region(QueryChainBufferSize, 0)

	ret := C.query_chain(unsafe.Pointer(regionReq), unsafe.Pointer(regionResponse))
	C.free(reqPtr)

	if ret != 0 {
		return nil, errors.New("failed to query chain")
	}

	response := TranslateToSlice(uintptr(regionResponse))
	return response, nil
}

// ------- query detail types ---------
type QueryResponse struct {
	Ok  []byte    `json:"Ok,omitempty"`
	Err *StdError `json:"Err,omitempty"`
}

// This is a 2-level result
type QuerierResult struct {
	Ok  *QueryResponse `json:"Ok,omitempty"`
	Err *SystemError   `json:"Err,omitempty"`
}

func ToQuerierResult(response []byte, err error) QuerierResult {
	if err == nil {
		return QuerierResult{
			Ok: &QueryResponse{
				Ok: response,
			},
		}
	}
	syserr := ToSystemError(err)
	if syserr != nil {
		return QuerierResult{
			Err: syserr,
		}
	}
	stderr := ToStdError(err)
	return QuerierResult{
		Ok: &QueryResponse{
			Err: stderr,
		},
	}
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
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

type StakingQuery struct {
	Validators     *ValidatorsQuery     `json:"validators,omitempty"`
	AllDelegations *AllDelegationsQuery `json:"all_delegations,omitempty"`
	Delegation     *DelegationQuery     `json:"delegation,omitempty"`
	BondedDenom    *struct{}            `json:"bonded_denom,omitempty"`
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
