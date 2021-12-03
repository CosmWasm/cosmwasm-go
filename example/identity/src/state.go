//go:generate ../../../bin/tinyjson -snake_case -all state.go
package src

// Person keeps track of the real identity of an address
//go:generate ../../../bin/generator state map state.go Person
type Person struct {
	Address    string `json:"address" orm:"primaryKey"`
	Name       string
	Surname    string
	City       string
	PostalCode int32
}
