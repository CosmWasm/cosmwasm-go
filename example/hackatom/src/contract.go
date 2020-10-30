package src

import (
	"errors"

	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// TODO: use querier
// TODO: add migration support
// TODO: last functions

func Init(deps *std.Deps, env std.Env, info std.MessageInfo, msg []byte) (*std.InitResultOk, error) {
	deps.Api.Debug("here we go 🚀")

	initMsg := InitMsg{}
	err := ezjson.Unmarshal(msg, &initMsg)
	if err != nil {
		return nil, err
	}

	// just verify these (later we save like that)
	_, err = deps.Api.CanonicalAddress(initMsg.Verifier)
	if err != nil {
		return nil, err
	}
	_, err = deps.Api.CanonicalAddress(initMsg.Beneficiary)
	if err != nil {
		return nil, err
	}

	state := State{
		Verifier:    initMsg.Verifier,
		Beneficiary: initMsg.Beneficiary,
		Funder:      info.Sender,
	}

	err = SaveState(deps.Storage, &state)
	if err != nil {
		return nil, err
	}
	res := std.InitResponse{
		Attributes: []std.Attribute{{"Let the", "hacking begin"}},
	}
	return &std.InitResultOk{Ok: res}, nil
}

func Handle(deps *std.Deps, env std.Env, info std.MessageInfo, data []byte) (*std.HandleResultOk, error) {
	msg := HandleMsg{}
	err := ezjson.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.Release.WasSet():
		return handleRelease(deps, &env, &info)
	case msg.CpuLoop.WasSet():
		return handleCpuLoop(deps, &env, &info)
	case msg.StorageLoop.WasSet():
		return handleStorageLoop(deps, &env, &info)
	case msg.MemoryLoop.WasSet():
		return handleMemoryLoop(deps, &env, &info)
	case msg.AllocateLargeMemory.WasSet():
		return nil, errors.New("Not implemented: AllocateLargeMemory")
	case msg.Panic.WasSet():
		return handlePanic(deps, &env, &info)
	case msg.UserErrorsInApiCalls.WasSet():
		return nil, errors.New("Not implemented: UserErrorInApiCalls")
	default:
		return nil, errors.New("Unknown HandleMsg")
	}
}

func handleRelease(deps *std.Deps, env *std.Env, info *std.MessageInfo) (*std.HandleResultOk, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if info.Sender != state.Verifier {
		return nil, errors.New("Unauthorized")
	}
	// TODO: pull this out into a helper
	query := std.QueryRequest{
		Bank: std.BankQuery{
			AllBalances: std.AllBalancesQuery{
				Address: env.Contract.Address,
			},
		},
	}
	deps.Api.Debug("ONE")
	binQuery, err := ezjson.Marshal(query)
	if err != nil {
		return nil, err
	}
	deps.Api.Debug("TWO")
	deps.Api.Debug(string(binQuery))
	data, err := deps.Querier.RawQuery(binQuery)
	if err != nil {
		return nil, err
	}
	deps.Api.Debug("THREE")
	deps.Api.Debug(string(data))
	var qres std.AllBalancesResponse
	err = ezjson.Unmarshal(data, &qres)
	if err != nil {
		return nil, err
	}
	deps.Api.Debug("FOUR")

	msg := []std.CosmosMsg{{
		Bank: std.BankMsg{
			Send: std.SendMsg{
				FromAddress: env.Contract.Address,
				ToAddress:   state.Beneficiary,
				Amount:      qres.Amount,
			},
		},
	}}

	res := std.HandleResponse{
		Attributes: []std.Attribute{
			{"action", "release"},
			{"destination", state.Beneficiary},
		},
		Messages: msg,
	}
	return &std.HandleResultOk{Ok: res}, nil
}

func handleCpuLoop(deps *std.Deps, env *std.Env, info *std.MessageInfo) (*std.HandleResultOk, error) {
	var counter uint64 = 0
	for {
		counter += 1
		if counter >= 9_000_000_000 {
			counter = 0
		}
	}
	return &std.HandleResultOk{}, nil
}

func handleMemoryLoop(deps *std.Deps, env *std.Env, info *std.MessageInfo) (*std.HandleResultOk, error) {
	counter := 1
	data := []int{1}
	for {
		counter += 1
		data = append(data, counter)
	}
	return &std.HandleResultOk{}, nil
}

func handleStorageLoop(deps *std.Deps, env *std.Env, info *std.MessageInfo) (*std.HandleResultOk, error) {
	var counter uint64 = 0
	for {
		data := []byte{0, 0, 0, 0, 0, 0, byte(counter / 256), byte(counter % 256)}
		deps.Storage.Set([]byte("test.key"), data)
	}
	return &std.HandleResultOk{}, nil
}

func handlePanic(deps *std.Deps, env *std.Env, info *std.MessageInfo) (*std.HandleResultOk, error) {
	panic("This page intentionally faulted")
}

func Query(deps *std.Deps, env std.Env, data []byte) (*std.QueryResponseOk, error) {
	msg := QueryMsg{}
	err := ezjson.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.Verifier.WasSet():
		return queryVerifier(deps, &env)
	case msg.OtherBalance.Address != "":
		return nil, errors.New("Not implemented: OtherBalance")
	case msg.Recurse.Work != 0 || msg.Recurse.Depth != 0:
		return nil, errors.New("Not implemented: Recurse")
	default:
		return nil, errors.New("Unknown QueryMsg")
	}
}

func queryVerifier(deps *std.Deps, env *std.Env) (*std.QueryResponseOk, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	res := VerifierResponse{
		Verifier: state.Verifier,
	}
	bz, err := ezjson.Marshal(res)
	if err != nil {
		return nil, err
	}

	return &std.QueryResponseOk{Ok: bz}, nil
}
