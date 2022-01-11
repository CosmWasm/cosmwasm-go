package src

import (
	imp "github.com/cosmwasm/cosmwasm-go/example/identityv2/src/imp"
)

type QueryMsg struct {
	QueryIdentity *QueryIdentity     `json:"query_identity"`
	QueryImported *imp.ImportedQuery `json:"query_imported"`
}

type ExecuteMsg struct {
	CreateIdentity  *MsgCreateIdentity   `json:"create_identity"`
	ImportedMessage *imp.ImportedMessage `json:"imported_message"`
}
