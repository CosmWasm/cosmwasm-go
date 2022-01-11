package src

// +cw:json
type MsgCreateIdentity struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	City       string `json:"city"`
	PostalCode int32  `json:"postal_code"`
}

// +cw:json
type MsgUpdateCity struct {
	City       string `json:"city"`
	PostalCode int32  `json:"postal_code"`
}

// +cw:json
type MsgDelete struct {
}

// +cw:json
type QueryIdentity struct {
	ID string `json:"id"`
}

type QueryIdentityResponse struct {
	Person *Person `json:"person"`
}

// +cw:json
type MsgMigrate struct {
}

// +cw:json
type MsgInstantiate struct {
}
