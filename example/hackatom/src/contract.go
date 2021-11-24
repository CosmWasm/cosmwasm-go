package src

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, msg []byte) (*types.Response, error) {
	deps.Api.Debug("here we go 🚀")

	initMsg := InitMsg{}
	err := initMsg.UnmarshalJSON(msg)
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
	res := &types.Response{
		Attributes: []types.EventAttribute{{"Let the", "hacking begin"}},
	}
	return res, nil
}

func Migrate(deps *std.Deps, env types.Env, msg []byte) (*types.Response, error) {
	migrateMsg := MigrateMsg{}
	err := migrateMsg.UnmarshalJSON(msg)
	if err != nil {
		return nil, err
	}

	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}
	state.Verifier = migrateMsg.Verifier
	err = SaveState(deps.Storage, state)
	if err != nil {
		return nil, err
	}

	res := &types.Response{Data: []byte("migrated")}
	return res, nil
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	msg := HandleMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.Release != nil:
		return executeRelease(deps, &env, &info)
	case msg.CpuLoop != nil:
		return executeCpuLoop(deps, &env, &info)
	case msg.StorageLoop != nil:
		return executeStorageLoop(deps, &env, &info)
	case msg.MemoryLoop != nil:
		return executeMemoryLoop(deps, &env, &info)
	case msg.AllocateLargeMemory != nil:
		return nil, types.GenericError("Not implemented: AllocateLargeMemory")
	case msg.Panic != nil:
		return executePanic(deps, &env, &info)
	case msg.UserErrorsInApiCalls != nil:
		return executeUserErrorsInApiCall(deps)
	default:
		return nil, types.GenericError("Unknown HandleMsg")
	}
}

// used to signal an empty error, since we cannot use fmt, nil errors would panic when added to a string
var noError = errors.New("<nil>")

// isGenericError is the function used to check comparison, we don't
// use errors.Is because types.GenericError is comparable and hence
// the package would check for the two errors to be identical.
// TODO: match MockApi error messages with the VM provided API error messages to use errors.Is
func isGenericError(err error) bool {
	_, ok := err.(types.GenericErr)
	return ok
}

func executeUserErrorsInApiCall(deps *std.Deps) (*types.Response, error) {
	// canonicalization

	// case empty
	_, err := deps.Api.CanonicalAddress("")
	if !isGenericError(err) {
		if err == nil {
			err = noError
		}
		return nil, errors.New("canonical empty unexpected error: " + err.Error())
	}
	// invalid bech32 addr
	_, err = deps.Api.CanonicalAddress("bn9hhssomeltvhzgvuqkwjkpwxojfuigltwedayzxljucefikuieillowaticksoistqoynmgcnj219a")
	if !isGenericError(err) {
		if err == nil {
			err = noError
		}
		return nil, errors.New("canonical invalid bech32 unexpected error: " + err.Error())
	}

	// humanization

	// empty
	_, err = deps.Api.HumanAddress([]byte{})
	if !isGenericError(err) {
		if err == nil {
			err = noError
		}
		return nil, errors.New("humanize empty unexpected error: " + err.Error())
	}

	// too short
	_, err = deps.Api.HumanAddress([]byte{0xAA, 0xBB, 0xCC})
	if !isGenericError(err) {
		if err == nil {
			err = noError
		}
		return nil, errors.New("humanize too short unexpected error: " + err.Error())
	}

	// wrong length
	_, err = deps.Api.HumanAddress([]byte{0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6})
	if !isGenericError(err) {
		if err == nil {
			err = noError
		}
		return nil, errors.New("humanize wrong length unexpected error " + err.Error())
	}

	return &types.Response{}, nil
}

func executeRelease(deps *std.Deps, env *types.Env, info *types.MessageInfo) (*types.Response, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if info.Sender != state.Verifier {
		return nil, types.Unauthorized{}
	}
	amount, err := std.QuerierWrapper{deps.Querier}.QueryAllBalances(env.Contract.Address)
	if err != nil {
		return nil, err
	}

	msg := types.NewSubMsg(types.CosmosMsg{
		Bank: &types.BankMsg{
			Send: &types.SendMsg{
				ToAddress: state.Beneficiary,
				Amount:    amount,
			},
		},
	})

	res := &types.Response{
		Attributes: []types.EventAttribute{
			{"action", "release"},
			{"destination", state.Beneficiary},
		},
		Messages: []types.SubMsg{msg},
	}
	return res, nil
}

func executeCpuLoop(deps *std.Deps, env *types.Env, info *types.MessageInfo) (*types.Response, error) {
	var counter uint64 = 0
	for {
		counter += 1
		if counter >= 9_000_000_000 {
			counter = 0
		}
	}
	return &types.Response{}, nil
}

func executeMemoryLoop(deps *std.Deps, env *types.Env, info *types.MessageInfo) (*types.Response, error) {
	counter := 1
	data := []int{1}
	for {
		counter += 1
		data = append(data, counter)
	}
	return &types.Response{}, nil
}

func executeStorageLoop(deps *std.Deps, env *types.Env, info *types.MessageInfo) (*types.Response, error) {
	var counter uint64 = 0
	for {
		data := []byte{0, 0, 0, 0, 0, 0, byte(counter / 256), byte(counter % 256)}
		deps.Storage.Set([]byte("test.key"), data)
	}
	return &types.Response{}, nil
}

func executePanic(deps *std.Deps, env *types.Env, info *types.MessageInfo) (*types.Response, error) {
	panic("This page intentionally faulted")
}

func Query(deps *std.Deps, env types.Env, data []byte) ([]byte, error) {
	msg := QueryMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	var res std.JSONType
	switch {
	case msg.Verifier != nil:
		res, err = queryVerifier(deps, &env)
	case msg.OtherBalance != nil:
		res, err = queryOtherBalance(deps, &env, msg.OtherBalance)
	case msg.Recurse != nil:
		return queryRecurse(deps, &env, msg.Recurse)
	default:
		err = types.GenericError("Unknown QueryMsg " + string(data))
	}
	if err != nil {
		return nil, err
	}

	// if we got a result above, encode it
	bz, err := res.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bz, nil

}

func queryRecurse(deps *std.Deps, env *types.Env, recurse *Recurse) ([]byte, error) {
	contractAddrBytes := []byte(env.Contract.Address)

	// perform work
	var result [32]byte
	for i := uint32(0); i < recurse.Work; i++ {
		result = sha256.Sum256(contractAddrBytes)
	}

	if recurse.Depth == 0 {
		return (RecurseResponse{
			Hashed: hex.EncodeToString(result[:]),
		}).MarshalJSON()
	}

	recurseRequest := Recurse{
		Depth: recurse.Depth - 1,
		Work:  recurse.Work,
	}

	recurseBytes, err := recurseRequest.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req := types.QueryRequest{
		Wasm: &types.WasmQuery{
			Smart: &types.SmartQuery{
				ContractAddr: env.Contract.Address,
				Msg:          recurseBytes,
			},
		},
	}

	reqBytes, err := req.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return deps.Querier.RawQuery(reqBytes)
}

func queryVerifier(deps *std.Deps, env *types.Env) (*VerifierResponse, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	return &VerifierResponse{
		Verifier: state.Verifier,
	}, nil
}

func queryOtherBalance(deps *std.Deps, env *types.Env, msg *OtherBalance) (*types.AllBalancesResponse, error) {
	amount, err := std.QuerierWrapper{Querier: deps.Querier}.QueryAllBalances(msg.Address)
	if err != nil {
		return nil, err
	}

	return &types.AllBalancesResponse{
		Amount: amount,
	}, nil
}
