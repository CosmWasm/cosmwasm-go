package src

import (
	"github.com/cosmwasm/cosmwasm-go/std"
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
	// TODO
	// data, err := storage.Get(StateKey)
	// if err != nil {
	// 	return nil, err
	// }

	var state State
	// err = state.UnmarshalJSON(data)
	// if err != nil {
	// 	return nil, err
	// }
	return &state, nil
}

func SaveState(storage std.Storage, state *State) error {
	// TODO
	return nil
	// bz, err := state.MarshalJSON()
	// if err != nil {
	// 	return err
	// }
	// // TODO: this should not return error
	// // TODO: change names Api -> Api
	// return storage.Set(StateKey, bz)
}
