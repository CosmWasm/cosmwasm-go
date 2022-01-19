package src

// MsgCreateIdentity creates a new Person
// +cw:json
type MsgCreateIdentity struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	City       string `json:"city"`
	PostalCode int32  `json:"postal_code"`
}

// MsgUpdateCity updates a Person city and postal code information
// +cw:json
type MsgUpdateCity struct {
	City       string `json:"city"`
	PostalCode int32  `json:"postal_code"`
}

// MsgDelete deletes the Person identity associated with the sender.
// +cw:json
type MsgDelete struct {
}

// +cw:json
type QueryIdentity struct {
	ID string `json:"id"`
}

// +cw:json
type QueryIdentityResponse struct {
	Person *Person `json:"person"`
}

// +cw:json
type MsgMigrate struct {
}

// +cw:json
type MsgInstantiate struct {
}
