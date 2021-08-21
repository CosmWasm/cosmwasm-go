module github.com/cosmwasm/cosmwasm-go

go 1.14

require (
	github.com/CosmWasm/go-cosmwasm v0.11.1-0.20201030003026-823a52c049ee
	github.com/cosmwasm/jsonparser v1.0.2
	github.com/mailru/easyjson v0.7.8
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tm-db v0.6.2
)

replace github.com/mailru/easyjson => github.com/CosmWasm/tinyjson v0.8.1-0.20210821135926-d2e906c18c4b