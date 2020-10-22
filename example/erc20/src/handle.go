package src

import (
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

func handleInvokeMessage(deps *std.Extern, env std.Env, msg []byte) (*std.HandleResultOk, *std.CosmosResponseError) {
	handler := Handler{}
	e := ezjson.Unmarshal(msg, &handler)
	if e != nil {
		return nil, std.GenerateError(std.GenericError, "Unamrshal handle error : "+e.Error(), "")
	}
	i := getNotEmptyElem(handler)
	state, e := LoadState(deps)
	if e != nil {
		return nil, std.GenerateError(std.GenericError, "LoadState error : "+e.Error(), "")
	}
	erc20 := NewErc20Protocol(state, deps, &env)
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
		return handleTransferOwner(deps, &env, i.(TransferOwner), ownerShip)
	case AcceptTransferredOwner:
		ownerShip.LoadOwner()
		return handleTransferOwnerAccepted(deps, &env, i.(AcceptTransferredOwner), ownerShip)
	default:
		return nil, std.GenerateError(std.GenericError, "Unsupported invoke type", "")
	}
}

func handleApprove(a Approve, erc20 Erc20) (*std.HandleResultOk, *std.CosmosResponseError) {
	if erc20.Approve([]byte(a.Spender), a.Value) {
		return std.HandleResultOkDefault(), nil
	}
	return nil, std.GenerateError(std.GenericError, "Approve failed", "")
}

func handleTransfer(t Transfer, erc20 Erc20) (*std.HandleResultOk, *std.CosmosResponseError) {
	if erc20.Transfer([]byte(t.ToAddr), t.Value) {
		return std.HandleResultOkDefault(), nil
	}
	return nil, std.GenerateError(std.GenericError, "Transfer failed", "")
}

func handleTransferFrom(tf TransferFrom, erc20 Erc20) (*std.HandleResultOk, *std.CosmosResponseError) {
	if erc20.TransferFrom([]byte(tf.FromAddr), []byte(tf.ToAddr), tf.Value) {
		return std.HandleResultOkDefault(), nil
	}
	return nil, std.GenerateError(std.GenericError, "TransferFrom failed", "")
}

func handleTransferOwner(deps *std.Extern, env *std.Env, to TransferOwner, owner Owner) (*std.HandleResultOk, *std.CosmosResponseError) {
	sender, err := deps.EApi.CanonicalAddress(env.Message.Sender)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Invalid Sender: "+err.Error(), "")
	}
	newOwner, err := deps.EApi.CanonicalAddress(to.NewOwner)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Invalid new owner: "+err.Error(), "")
	}

	owner.TransferOwnership(sender, newOwner)
	if owner.SaveOwner() {
		return std.HandleResultOkDefault(), nil
	}
	return nil, std.GenerateError(std.GenericError, "Transfer ownership execute success", "")
}

func handleTransferOwnerAccepted(deps *std.Extern, env *std.Env, atf AcceptTransferredOwner, owner Owner) (*std.HandleResultOk, *std.CosmosResponseError) {
	sender, err := deps.EApi.CanonicalAddress(env.Message.Sender)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Invalid Sender: "+err.Error(), "")
	}
	accepted, err := deps.EApi.CanonicalAddress(atf.AcceptedAddr)
	if err != nil {
		return nil, std.GenerateError(std.GenericError, "Invalid accetped owner: "+err.Error(), "")
	}

	owner.AcceptTransfer(sender, accepted)
	if owner.SaveOwner() {
		return std.HandleResultOkDefault(), nil
	}
	return nil, std.GenerateError(std.GenericError, "Transfer ownership execute success", "")
}

//query handler

func getNotEmptyQueryElem(q Querier) interface{} {
	if len(q.Balance.Address) > 0 {
		return q.Balance
	}

	return nil
}

func handleQuery(deps *std.Extern, msg []byte) (*std.QueryResponseOk, *std.CosmosResponseError) {
	querier := Querier{}
	e := ezjson.Unmarshal(msg, &querier)
	if e != nil {
		return nil, std.GenerateError(std.GenericError, "Unamrshal handle error : "+e.Error(), "")
	}
	state, e := LoadState(deps)
	if e != nil {
		return nil, std.GenerateError(std.GenericError, "LoadState error : "+e.Error(), "")
	}
	erc20 := NewErc20Protocol(state, deps, nil)
	q := getNotEmptyQueryElem(querier)

	switch q.(type) {
	case Balance:
		return QueryBalance(q.(Balance), erc20)
	}
	return nil, std.GenerateError(std.GenericError, "Unsupported query type", "")
}

func QueryBalance(b Balance, erc20 Erc20) (*std.QueryResponseOk, *std.CosmosResponseError) {
	br := BalanceResponse{Value: erc20.BalanceOf(b.Address)}
	v, e := ezjson.Marshal(br)
	if e != nil {
		return nil, std.GenerateError(std.GenericError, "Marshal BalanceResponse Failed "+e.Error(), "")
	}
	return std.BuildQueryResponse(string(v)), nil
}
