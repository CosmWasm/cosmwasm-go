module github.com/cosmwasm/cosmwasm-go

go 1.14

require (
	github.com/CosmWasm/wasmvm v0.13.1
	github.com/mailru/easyjson v0.7.8
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tm-db v0.6.2
)

replace github.com/mailru/easyjson => github.com/CosmWasm/tinyjson v0.8.2
