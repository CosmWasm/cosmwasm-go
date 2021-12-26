//go:generate ../../../bin/tinyjson -snake_case -all msgs.go
package src

type MsgCreateIdentity struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	City       string `json:"city"`
	PostalCode int32  `json:"postal_code"`
}

type MsgUpdateCity struct {
	City       string `json:"city"`
	PostalCode int32  `json:"postal_code"`
}

type MsgDelete struct{}

type QueryIdentity struct {
	ID string `json:"id"`
}
