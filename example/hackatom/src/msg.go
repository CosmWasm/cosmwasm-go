package src

import "github.com/cosmwasm/cosmwasm-go/std/ezjson"

type InitMsg struct {
	Verifier    string `json:"verifier"`
	Beneficiary string `json:"beneficiary"`
}

type MigrateMsg struct {
	Verifier string `json:"verifier"`
}

type HandleMsg struct {
	Release              ezjson.EmptyStruct `json:"release,opt_seen"`
	CpuLoop              ezjson.EmptyStruct `json:"cpu_loop,opt_seen"`
	StorageLoop          ezjson.EmptyStruct `json:"storage_loop,opt_seen"`
	MemoryLoop           ezjson.EmptyStruct `json:"memory_loop,opt_seen"`
	AllocateLargeMemory  ezjson.EmptyStruct `json:"allocate_large_memory,opt_seen"`
	Panic                ezjson.EmptyStruct `json:"panic,opt_seen"`
	UserErrorsInApiCalls ezjson.EmptyStruct `json:"user_errors_in_api_calls,opt_seen"`
}

type QueryMsg struct {
	Verifier     ezjson.EmptyStruct `json:"verifier,opt_seen"`
	OtherBalance OtherBalance       `json:"other_balance,omit_empty"`
	Recurse      Recurse            `json:"recurse,omit_empty"`
}

type OtherBalance struct {
	Address string `json:"address"`
}

type Recurse struct {
	Depth uint32 `json:"depth"`
	Work  uint32 `json:"work"`
}

type VerifierResponse struct {
	Verifier string `json:"verifier"`
}

type RecurseResponse struct {
	// this should be base64 binary - we just encode it manually outside of ezjson
	Hashed string `json:"hashed"`
}
