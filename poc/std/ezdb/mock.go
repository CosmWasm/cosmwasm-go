// +build !cosmwasm

package ezdb

import "errors"

var storage map[string][]byte

func init() {
	storage = make(map[string][]byte)
}

func WriteStorage(key []byte,value []byte) (error){
	storage[string(key)] = value
	return nil
}

func ReadStorage(key []byte) ([]byte,error){
	value, ok := storage[string(key)]
	// TODO: revisit if we return nil or error in this case
	if !ok {
		return nil, errors.New("No data found")
	}
	return value, nil
}
