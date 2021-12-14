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

type MsgDelete struct {
}

// ExecuteMsg is used to execute state transitions on the identity contract
type ExecuteMsg struct {
	CreateIdentity *MsgCreateIdentity `json:"create_identity"`
	UpdateCity     *MsgUpdateCity     `json:"update_city"`
	DeleteIdentity *MsgDelete         `json:"delete_identity"`
}

type QueryIdentity struct {
	ID string `json:"id"`
}

type QueryMsg struct {
	Identity *QueryIdentity `json:"identity"`
}
