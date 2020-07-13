package src

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
