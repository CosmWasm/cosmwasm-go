package src

import (
	"strconv"

	"github.com/CosmWasm/cosmwasm-go/example/voter/src/state"
	"github.com/CosmWasm/cosmwasm-go/example/voter/src/types"
	"github.com/CosmWasm/cosmwasm-go/std"
	stdTypes "github.com/CosmWasm/cosmwasm-go/std/types"
)

// Instantiate performs the contract state initialization.
func Instantiate(deps *std.Deps, env stdTypes.Env, info stdTypes.MessageInfo, msgBz []byte) (*stdTypes.Response, error) {
	deps.Api.Debug("Instantiate called")

	var msg types.MsgInstantiate
	if err := msg.UnmarshalJSON(msgBz); err != nil {
		return nil, types.NewErrInvalidRequest("msg JSON unmarshal: " + err.Error())
	}

	return handleMsgInstantiate(deps, info, msg)
}

// Migrate performs the contract state upgrade and can only be called by the contract admin.
// If the admin field is not set for a contract, contract is immutable.
func Migrate(deps *std.Deps, env stdTypes.Env, msgBz []byte) (*stdTypes.Response, error) {
	return nil, types.NewErrUnimplemented("Migrate")
}

// Execute performs the contract state change.
func Execute(deps *std.Deps, env stdTypes.Env, info stdTypes.MessageInfo, msgBz []byte) (*stdTypes.Response, error) {
	deps.Api.Debug("Execute called")

	var msg types.MsgExecute
	if err := msg.UnmarshalJSON(msgBz); err != nil {
		return nil, types.NewErrInvalidRequest("msg JSON unmarshal: " + err.Error())
	}

	switch {
	case msg.Release != nil:
		return handleMsgRelease(deps, env, info)
	case msg.NewVoting != nil:
		return handleMsgNewVoting(deps, env, info, *msg.NewVoting)
	case msg.Vote != nil:
		return handleMsgVote(deps, env, info, *msg.Vote)
	}

	return nil, types.NewErrInvalidRequest("unknown execute request")
}

// Sudo performs the contract state change and can only be called by a native Cosmos module (like x/gov).
func Sudo(deps *std.Deps, env stdTypes.Env, msgBz []byte) (*stdTypes.Response, error) {
	deps.Api.Debug("Sudo called")

	var msg types.MsgSudo
	if err := msg.UnmarshalJSON(msgBz); err != nil {
		return nil, types.NewErrInvalidRequest("msg JSON unmarshal: " + err.Error())
	}

	switch {
	case msg.ChangeNewVotingCost != nil:
		return handleSudoChangeNewVotingCost(deps, *msg.ChangeNewVotingCost)
	case msg.ChangeVoteCost != nil:
		return handleSudoChangeVoteCost(deps, *msg.ChangeVoteCost)
	}

	return nil, types.NewErrInvalidRequest("unknown sudo request")
}

// Reply performs an optional contract state change.
// Endpoint is called when stdTypes.SubMsg was sent with (always/success/error) ReplyOn policy
// on other endpoint invocation (instantiate/execute/migrate/sudo/reply).
// SubMsg identification is done via stdTypes.SubMsg.ID field.
func Reply(deps *std.Deps, env stdTypes.Env, reply stdTypes.Reply) (*stdTypes.Response, error) {
	deps.Api.Debug("Reply called")

	replyType, found, err := state.GetReplyMsgType(deps.Storage, reply.ID)
	if err != nil {
		return nil, types.NewErrInternal(err.Error())
	}
	if !found {
		return nil, types.NewErrInternal("replyID (" + strconv.FormatUint(reply.ID, 10) + "): not found")
	}

	switch replyType {
	case state.ReplyMsgTypeBank:
		return handleReplyBankMsg(deps, reply.Result)
	}

	return nil, types.NewErrInternal("unknown replyMsgType: " + strconv.Itoa(int(replyType)))
}

// Query performs the contract state read.
func Query(deps *std.Deps, env stdTypes.Env, msgBz []byte) ([]byte, error) {
	deps.Api.Debug("Query called")

	var msg types.MsgQuery
	if err := msg.UnmarshalJSON(msgBz); err != nil {
		return nil, types.NewErrInvalidRequest("msg JSON unmarshal: " + err.Error())
	}

	var handlerRes std.JSONType
	var handlerErr error
	switch {
	case msg.Params != nil:
		handlerRes, handlerErr = queryParams(deps)
	case msg.Voting != nil:
		handlerRes, handlerErr = queryVoting(deps, *msg.Voting)
	case msg.Tally != nil:
		handlerRes, handlerErr = queryTally(deps, env, *msg.Tally)
	case msg.Open != nil:
		handlerRes, handlerErr = queryOpen(deps, env)
	case msg.ReleaseStats != nil:
		handlerRes, handlerErr = queryReleaseStats(deps)
	default:
		handlerErr = types.NewErrInvalidRequest("unknown query")
	}
	if handlerErr != nil {
		return nil, handlerErr
	}

	resBz, err := handlerRes.MarshalJSON()
	if err != nil {
		return nil, types.NewErrInternal("query result JSON marshal: " + err.Error())
	}

	return resBz, nil
}
