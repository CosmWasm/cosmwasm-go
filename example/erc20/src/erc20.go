package src

import (
	"bytes"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
	"github.com/cosmwasm/cosmwasm-go/std/safe_math"
)

type Owner interface {
	//init owner once
	Owned(ownerAddr []byte)

	//get current owner address
	GetOwner() []byte

	//get new owner address
	GetNewOwner() []byte

	//transfer contract owner from `sender` to `newer`
	TransferOwnership(sender, newer []byte)

	//check is `addr` is an ownership of this contract
	OnlyOwner(addr []byte) bool

	//event of ownership transferred
	OwnershipTransferred(from, to []byte)

	AcceptTransfer(sender, to []byte)

	SaveOwner() bool
	LoadOwner() bool
}

type Erc20 interface {
	Name() string
	Symbol() string
	Decimals() uint64
	TotalSupply() uint64
	BalanceOf(addr []byte) uint64
	Transfer(toAddr []byte, value uint64) bool
	TransferFrom(form, to []byte, value uint64) bool
	Approve(spender []byte, value uint64) bool
	EventOfTransfer(from, to []byte, value uint64)
	EventOfApproval(owner, spender []byte, value uint64)
	Assign(addr []byte, value uint64)
	//Management of contract meta data
	SaveState() bool
}

type Ownership struct {
	owner    []byte
	newOwner []byte

	apis *std.Extern
}

func (o *Ownership) Owned(ownerAddr []byte) {
	if len(o.owner) > 0 { //only once
		return
	}
	o.owner = ownerAddr
}

func (o Ownership) GetOwner() []byte {
	return o.owner
}

func (o Ownership) GetNewOwner() []byte {
	return o.newOwner
}

func (o *Ownership) TransferOwnership(sender, newer []byte) {
	if !o.OnlyOwner(sender) {
		//from must an owner of this contract
		return
	}
	o.newOwner = newer
}

func (o Ownership) OnlyOwner(addr []byte) bool {
	return bytes.Equal(o.owner, addr)
}

func (o Ownership) OwnershipTransferred(sender, to []byte) {
}

func (o *Ownership) AcceptTransfer(sender, to []byte) {
	if !bytes.Equal(sender, o.newOwner) || !bytes.Equal(to, o.newOwner) {
		return
	}
	o.owner = o.newOwner
	o.newOwner = []byte("")
}

func (o Ownership) SaveOwner() bool {
	//unhandled error
	if o.owner != nil {
		o.apis.EStorage.Set([]byte("owner"), o.owner)
	}

	if o.newOwner != nil {
		o.apis.EStorage.Set([]byte("newOwner"), o.newOwner)
	}

	return true
}

func (o *Ownership) LoadOwner() bool {
	//unhandled error
	o.owner, _ = o.apis.EStorage.Get([]byte("owner"))
	o.newOwner, _ = o.apis.EStorage.Get([]byte("newOwner"))

	return true
}

type State struct {
	NameOfToken   string `json:"name"`
	SymbolOfToken string `json:"symbol"`
	DecOfTokens   uint64 `json:"decimals"`
	TotalSupplyOf uint64 `json:"total_supply"`
}

type implErc20 struct {
	State
	apis *std.Extern
	env  *std.Env
}

func approvalPrefix(addr []byte) []byte {
	return []byte("approval" + string(addr))
}

func amountPrefix(addr []byte) []byte {
	return []byte("amount" + string(addr))
}

func (i implErc20) Name() string {
	return i.NameOfToken
}

func (i implErc20) Symbol() string {
	return i.SymbolOfToken
}

func (i implErc20) Decimals() uint64 {
	return i.DecOfTokens
}

func (i implErc20) TotalSupply() uint64 {
	return i.TotalSupplyOf
}

func (i implErc20) BalanceOf(addr []byte) uint64 {
	v, e := i.apis.EStorage.Get(amountPrefix(addr))
	if e != nil {
		return 0
	}
	return std.BytesToUint64(v)
}

func (i implErc20) getApproval(addr []byte) uint64 {
	v, ea := i.apis.EStorage.Get(approvalPrefix(addr))
	if ea != nil {
		return 0
	}
	return std.BytesToUint64(v)
}

func (i implErc20) setApproval(addr []byte, value uint64) bool {
	ea := i.apis.EStorage.Set(approvalPrefix(addr), std.Uint64toBytes(value))
	if ea != nil {
		return false
	}
	return true
}

func (i implErc20) Assign(addr []byte, value uint64) {
	i.apis.EStorage.Set(amountPrefix(addr), std.Uint64toBytes(value))
}

func (i implErc20) transfer(from, to []byte, value uint64) bool {
	m := i.BalanceOf(from)
	if m < value {
		return false
	}
	tm := i.BalanceOf(to)
	sender_money, es := safe_math.SafeSub(m, value)
	reciver_money, er := safe_math.SafeAdd(tm, value)
	if es != nil || er != nil {
		return false
	}

	es = i.apis.EStorage.Set(amountPrefix(from), std.Uint64toBytes(sender_money))
	er = i.apis.EStorage.Set(amountPrefix(to), std.Uint64toBytes(reciver_money))

	m = i.BalanceOf(to)
	if es != nil || er != nil {
		return false
	}
	return true
}

func (i implErc20) Transfer(toAddr []byte, value uint64) bool {
	sender, err := i.apis.EApi.CanonicalAddress(i.env.Message.Sender)
	if err != nil {
		// TODO: use an error
		//return nil, std.GenerateError(std.GenericError, "Invalid Sender: " + err.Error(), "")
		panic("invalid sender - expected valid bech32")
	}
	return i.transfer(sender, toAddr, value)
}

func (i implErc20) TransferFrom(from, to []byte, value uint64) bool {
	approval := i.getApproval(from)
	if approval == 0 || value > approval {
		return false
	}
	return i.transfer(from, to, value)
}

func (i implErc20) Approve(spender []byte, value uint64) bool {
	m := i.BalanceOf(spender)
	if m < value {
		return false
	}
	return i.setApproval(spender, value)
}

func (i implErc20) EventOfTransfer(from, to []byte, value uint64) {
	return
}

func (i implErc20) EventOfApproval(owner, spender []byte, value uint64) {
	return
}

func (i implErc20) SaveState() bool {
	b, e := ezjson.Marshal(i.State)
	if e != nil {
		return false
	}
	e = i.apis.EStorage.Set([]byte("State"), b)
	return e == nil
}

func LoadState(extern *std.Extern) (State, error) {
	state := State{}
	v, e := extern.EStorage.Get([]byte("State"))
	if e != nil {
		return state, e
	}
	e = ezjson.Unmarshal(v, &state)
	return state, e
}

func NewErc20Protocol(state State, extern *std.Extern, env *std.Env) Erc20 {
	return implErc20{
		State: state,
		apis:  extern,
		env:   env,
	}
}

func NewOwnership(extern *std.Extern) Owner {
	return &Ownership{owner: nil, newOwner: nil, apis: extern}
}
