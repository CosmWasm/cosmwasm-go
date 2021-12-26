//go:generate ../../../bin/tinyjson -snake_case -all state.go
package src

// Person keeps track of the real identity of an address
//go:generate ../../../bin/generator state map state.go Person
type Person struct {
	Address    string `json:"address,omitempty" orm:"primaryKey"`
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode int32  `json:"postal_code,omitempty"`
}
