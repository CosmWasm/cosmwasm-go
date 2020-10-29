package src

import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func getNotEmptyElem(h Handler) interface{} {
	if len(h.Approve.Spender) > 0 {
		return h.Approve
	}
	if len(h.Transfer.ToAddr) > 0 {
		return h.Transfer
	}
	if len(h.TransferFrom.FromAddr) > 0 && len(h.TransferFrom.ToAddr) > 0 {
		return h.TransferFrom
	}
	if len(h.TransferOwner.NewOwner) > 0 {
		return h.TransferOwner
	}
	if len(h.AcceptedTransferOwner.AcceptedAddr) > 0 {
		return h.AcceptedTransferOwner
	}
	return nil
}

func handleInvokeMessage(deps *std.Extern, env std.Env, info std.MessageInfo, msg []byte) (*std.HandleResultOk, error) {
	handler := Handler{}
	err := ezjson.Unmarshal(msg, &handler)
	if err != nil {
		return nil, err
	}
	i := getNotEmptyElem(handler)
	state, err := LoadState(deps)
	if err != nil {
		return nil, err
	}
	erc20 := NewErc20Protocol(state, deps, &info)
	ownerShip := NewOwnership(deps)
	switch i.(type) {
	case Approve:
		return handleApprove(i.(Approve), erc20)
	case Transfer:
		return handleTransfer(i.(Transfer), erc20)
	case TransferFrom:
		return handleTransferFrom(i.(TransferFrom), erc20)
	case TransferOwner:
		ownerShip.LoadOwner()
		return handleTransferOwner(deps, &info, i.(TransferOwner), ownerShip)
	case AcceptTransferredOwner:
		ownerShip.LoadOwner()
		return handleTransferOwnerAccepted(deps, &info, i.(AcceptTransferredOwner), ownerShip)
	default:
		return nil, errors.New("Unsupported HandleMsg variant")
	}
}

func handleApprove(a Approve, erc20 Erc20) (*std.HandleResultOk, error) {
	if erc20.Approve([]byte(a.Spender), a.Value) {
		return std.HandleResultOkDefault(), nil
	}
	return nil, errors.New("Approve failed")
}

func handleTransfer(t Transfer, erc20 Erc20) (*std.HandleResultOk, error) {
	if erc20.Transfer([]byte(t.ToAddr), t.Value) {
		return std.HandleResultOkDefault(), nil
	}
	return nil, errors.New("Transfer failed")
}

func handleTransferFrom(tf TransferFrom, erc20 Erc20) (*std.HandleResultOk, error) {
	if erc20.TransferFrom([]byte(tf.FromAddr), []byte(tf.ToAddr), tf.Value) {
		return std.HandleResultOkDefault(), nil
	}
	return nil, errors.New("TransferFrom failed")
}

func handleTransferOwner(deps *std.Extern, info *std.MessageInfo, to TransferOwner, owner Owner) (*std.HandleResultOk, error) {
	sender, err := deps.Api.CanonicalAddress(info.Sender)
	if err != nil {
		return nil, err
	}
	newOwner, err := deps.Api.CanonicalAddress(to.NewOwner)
	if err != nil {
		return nil, err
	}

	owner.TransferOwnership(sender, newOwner)
	if owner.SaveOwner() {
		return std.HandleResultOkDefault(), nil
	}
	return nil, errors.New("TransferOwner save failed")
}

func handleTransferOwnerAccepted(deps *std.Extern, info *std.MessageInfo, atf AcceptTransferredOwner, owner Owner) (*std.HandleResultOk, error) {
	sender, err := deps.Api.CanonicalAddress(info.Sender)
	if err != nil {
		return nil, err
	}
	accepted, err := deps.Api.CanonicalAddress(atf.AcceptedAddr)
	if err != nil {
		return nil, err
	}

	owner.AcceptTransfer(sender, accepted)
	if owner.SaveOwner() {
		return std.HandleResultOkDefault(), nil
	}
	return nil, errors.New("TransferOwnerAccepted save failed")
}

//query handler

func getNotEmptyQueryElem(q Querier) interface{} {
	if len(q.Balance.Address) > 0 {
		return q.Balance
	}

	return nil
}

func handleQuery(deps *std.Extern, _env std.Env, msg []byte) (*std.QueryResponseOk, error) {
	querier := Querier{}
	err := ezjson.Unmarshal(msg, &querier)
	if err != nil {
		return nil, err
	}
	state, err := LoadState(deps)
	if err != nil {
		return nil, err
	}
	erc20 := NewErc20Protocol(state, deps, nil)
	q := getNotEmptyQueryElem(querier)

	switch q.(type) {
	case Balance:
		return QueryBalance(q.(Balance), erc20)
	}
	return nil, errors.New("Unsupported query type")
}

func QueryBalance(b Balance, erc20 Erc20) (*std.QueryResponseOk, error) {
	br := BalanceResponse{Value: erc20.BalanceOf(b.Address)}
	v, err := ezjson.Marshal(br)
	if err != nil {
		return nil, err
	}
	return std.BuildQueryResponse(string(v)), nil
}
