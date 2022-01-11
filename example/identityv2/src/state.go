package src

// Person keeps track of the real identity of an address
// +cw:json
// +cw:state map
type Person struct {
	Address    string `json:"address" orm:"primaryKey"`
	Name       string
	Surname    string
	City       string
	PostalCode int32
}

type Bau string
