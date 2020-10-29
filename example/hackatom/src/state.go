package src

import (
	"fmt"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// this is what we store
type State struct {
	// having issues storing []byte (maybe human addr is better?)
	Owner string `json:"owner"`
	// Huh, why does this fail to marhsall as uint64?
	//Count uint64 `json:"count"`
	Count int64 `json:"count"`
}

var StateKey = []byte("State")

func LoadState(storage std.Storage) (*State, error) {
	var state State
	data, err := storage.Get(StateKey)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	err = ezjson.Unmarshal(data, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func SaveState(storage std.Storage, state *State) error {
	bz, err := ezjson.Marshal(*state)
	if err != nil {
		return err
	}
	// TODO: this should not return error
	// TODO: change names EApi -> Api
	return storage.Set(StateKey, bz)
}
