//go:generate ../../../bin/tinyjson -all -snake_case contract.go
package src

import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

// FirstKey defines the value of the default key,
// when no key is set in the contract so far.
// NOTE: keys are [1]byte length but in KV they're [n]bytes.
const FirstKey byte = 0

// Item defines the state object of the queue values.
type Item struct {
	Value int32 `json:"value"`
}

// Enqueue is the message used to add an Item to the queue.
type Enqueue struct {
	Value int32 `json:"value"`
}

// Dequeue is the message used to remove an Item from the queue.
type Dequeue struct{}

// ExecuteMsg defines all the messages that modify state that can be sent to the contract.
type ExecuteMsg struct {
	// Enqueue adds a value in the queue
	Enqueue *Enqueue `json:"enqueue"`
	// Dequeue removes a value from the queue
	Dequeue *Dequeue `json:"dequeue"`
}

// QueryMsg defines all the set of the possible queries that can be sent to the contract.
type QueryMsg struct {
	// Count counts how many items in the queue
	Count *struct{} `json:"count"`
	// Sum the number of values in the queue
	Sum *struct{} `json:"sum"`
	// Reducer keeps open two iters at once
	Reducer *struct{} `json:"reducer"`
	// List lists
	List *struct{} `json:"list"`
}

// InstantiateMsg is the instantiation messages.
type InstantiateMsg struct{}

// Instantiate does nothing.
func Instantiate(_ *std.Deps, _ types.Env, _ types.MessageInfo, _ []byte) (*types.Response, error) {
	return &types.Response{}, nil
}

// Execute runs state modifying handlers of the contract given msg data.
func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	msg := ExecuteMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	switch {
	case msg.Enqueue != nil:
		return executeEnqueue(deps, env, info, msg.Enqueue)
	case msg.Dequeue != nil:
		return executeDequeue(deps, env, info, msg.Dequeue)
	}
	return nil, errors.New("unknown request") // TODO(fdymylja): make this a common error in some package once we sort out devex on ExecuteMsg
}

func executeDequeue(deps *std.Deps, _ types.Env, _ types.MessageInfo, _ *Dequeue) (*types.Response, error) {
	iter := deps.Storage.Range(nil, nil, std.Ascending)
	resp := &types.Response{}
	k, v, err := iter.Next()
	// if queue is empty, then return an empty response to signal nothing was dequeued
	if err != nil {
		return resp, nil
	}
	// otherwise, delete the key and return the removed value
	deps.Storage.Remove(k)
	resp.Data = v
	return resp, nil
}

func executeEnqueue(deps *std.Deps, _ types.Env, _ types.MessageInfo, enqueue *Enqueue) (*types.Response, error) {
	iter := deps.Storage.Range(nil, nil, std.Descending)
	nextKey, _, err := iter.Next()
	switch err {
	// if no error then increase the key by 1
	case nil:
		nextKey[0] = nextKey[0] + 1
	// otherwise, set next key as FirstKey
	default:
		nextKey = []byte{FirstKey}
	}
	value, err := (Item{Value: enqueue.Value}).MarshalJSON()
	if err != nil {
		return nil, err
	}
	deps.Storage.Set(nextKey, value)
	return &types.Response{}, nil
}

// Migrate executes queue contract's migration which consists in clearing
// the state and writing three new values in the queue
func Migrate(deps *std.Deps, _ types.Env, _ []byte) (*types.Response, error) {
	iter := deps.Storage.Range(nil, nil, std.Ascending)
	// clear
	for k, _, err := iter.Next(); err != nil; {
		deps.Storage.Remove(k)
	}
	// add three values
	for i := int32(100); i < 103; i++ {
		_, err := executeEnqueue(deps, types.Env{}, types.MessageInfo{}, &Enqueue{Value: i})
		if err != nil {
			return nil, err
		}
	}

	return &types.Response{}, nil
}

func Query(deps *std.Deps, env types.Env, msg []byte) ([]byte, error) {
	panic("impl")
}
