package main

import (
	"github.com/cosmwasm/cosmwasm-go/std/ezjson"
)

func main() {}

//export add
func add(a uint32, b uint32) uint32 {
	return a + b
}

type Foo struct {
	Num int32 `json:"num"`
}

//export parse
func parse() int32 {
	// parse just runs the ezjson code on a minimal type to see if this imports floats
	msg := []byte(`{"num":1234`)
	foo := Foo{}
	e := ezjson.Unmarshal(msg, &foo)
	if e != nil {
		panic(e)
	}
	return foo.Num
}