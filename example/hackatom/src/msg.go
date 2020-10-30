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
	Release              ezjson.EmptyStruct `json:"release,opt_seen,omitempty"`
	CpuLoop              ezjson.EmptyStruct `json:"cpu_loop,opt_seen,omitempty"`
	StorageLoop          ezjson.EmptyStruct `json:"storage_loop,opt_seen,omitempty"`
	MemoryLoop           ezjson.EmptyStruct `json:"memory_loop,opt_seen,omitempty"`
	AllocateLargeMemory  ezjson.EmptyStruct `json:"allocate_large_memory,opt_seen,omitempty"`
	Panic                ezjson.EmptyStruct `json:"panic,opt_seen,omitempty"`
	UserErrorsInApiCalls ezjson.EmptyStruct `json:"user_errors_in_api_calls,opt_seen,omitempty"`
}

type QueryMsg struct {
	Verifier     ezjson.EmptyStruct `json:"verifier,opt_seen,omitempty"`
	OtherBalance OtherBalance       `json:"other_balance,omitempty"`
	Recurse      Recurse            `json:"recurse,omitempty"`
}

type OtherBalance struct {
	Address string `json:"address,omitempty"`
}

type Recurse struct {
	Depth uint32 `json:"depth,omitempty"`
	Work  uint32 `json:"work,omitempty"`
}

type VerifierResponse struct {
	Verifier string `json:"verifier"`
}

type RecurseResponse struct {
	// this should be base64 binary - we just encode it manually outside of ezjson
	Hashed string `json:"hashed"`
}
