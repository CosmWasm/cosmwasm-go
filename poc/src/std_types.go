package src

import (
	"encoding/base64"

	"github.com/cosmwasm/cosmwasm-go/poc/std/ezjson"
)

//**** all these types can move to the standard lib once they are cleaned up ****//

type OKResponse struct {
	Ok string
}

func newOkResponse(resp string) *OKResponse {
	return &OKResponse{
		Ok: resp,
	}
}

func (o OKResponse) ToJSON() []byte {
	msg := `{"Ok":` + o.Ok + `}`
	return []byte(msg)
}

type ERRResponse struct {
	Err string
}

func newErrResponse(resp string) *ERRResponse {
	return &ERRResponse{
		Err: resp,
	}
}

func (e ERRResponse) ToJSON() []byte {
	msg := `{"Err":` + e.Err + `}`
	return []byte(msg)
}

type InitResponse struct {
	Data     string
	Log      string
	Messages string
}

// this takes a positive response and tries to return an OkResponse object.
// if there are encoding errors, then return *ERRResponse
func EncodeResult(obj interface{}) (*OKResponse, *ERRResponse) {
	bz, err := ezjson.MarshalEx(obj)
	if err != nil {
		return nil, newErrResponse(err.Error())
	}
	// if OKResponse took bytes, we wouldn't need this
	encoded := `"` + base64.StdEncoding.EncodeToString(bz) + `"`
	return newOkResponse(encoded), nil
}
