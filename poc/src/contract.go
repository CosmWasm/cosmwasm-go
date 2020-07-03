package src

import (
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
)

type InitMsg struct{
	Verifier string
	Beneficiary string
	Legacy 		int
}

type OKResponse struct {
	Ok string
}

type ERRResponse struct {
	Err string
}

type InitResponse struct {
	Data string
	Log string
	Messages string
}

func newOkResponse(resp string) *OKResponse{
	return &OKResponse{
		Ok: resp,
	}
}

func newErrResponse(resp string) *ERRResponse{
	return &ERRResponse{
		Err: resp,
	}
}


func go_init(msg InitMsg)(*OKResponse,*ERRResponse){
	initResp := InitResponse{
		Data: "",
		Log: "failed",
		Messages: "nil message",
	}
	if msg.Verifier == msg.Beneficiary {
		return nil,newErrResponse("error")
	}
	initResp.Log = "success"
	b,e := ezjson.Marshal(initResp)
	if e != nil {
		return nil,newErrResponse("error")
	}
	return newOkResponse(string(b)),nil
}

func go_handle(){

}

func go_query(){

}