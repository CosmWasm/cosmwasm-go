package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

// this is what we store
type State struct {
	// TODO: convert to canonical addresses when that is supported by ezjson
	Verifier    string `json:"VERIFIER"`
	Beneficiary string `json:"BENEFICIARY"`
	Funder      string `json:"FUNDER"`
}

var StateKey = []byte("config")

func LoadState(storage std.Storage) (*State, error) {
	var state State
	data, err := storage.Get(StateKey)
	if err != nil {
		return nil, err
	}
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
	// TODO: change names Api -> Api
	return storage.Set(StateKey, bz)
}
