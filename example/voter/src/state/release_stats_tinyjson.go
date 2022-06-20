// Code generated by tinyjson for marshaling/unmarshaling. DO NOT EDIT.

package state

import (
	types "github.com/CosmWasm/cosmwasm-go/std/types"
	tinyjson "github.com/CosmWasm/tinyjson"
	jlexer "github.com/CosmWasm/tinyjson/jlexer"
	jwriter "github.com/CosmWasm/tinyjson/jwriter"
)

// suppress unused package warning
var (
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ tinyjson.Marshaler
)

func tinyjson16010b1DecodeGithubComCosmWasmCosmwasmGoExampleVoterSrcState(in *jlexer.Lexer, out *ReleaseStats) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "count":
			out.Count = uint64(in.Uint64())
		case "total_amount":
			if in.IsNull() {
				in.Skip()
				out.TotalAmount = nil
			} else {
				in.Delim('[')
				if out.TotalAmount == nil {
					if !in.IsDelim(']') {
						out.TotalAmount = make([]types.Coin, 0, 2)
					} else {
						out.TotalAmount = []types.Coin{}
					}
				} else {
					out.TotalAmount = (out.TotalAmount)[:0]
				}
				for !in.IsDelim(']') {
					var v1 types.Coin
					(v1).UnmarshalTinyJSON(in)
					out.TotalAmount = append(out.TotalAmount, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson16010b1EncodeGithubComCosmWasmCosmwasmGoExampleVoterSrcState(out *jwriter.Writer, in ReleaseStats) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Count))
	}
	{
		const prefix string = ",\"total_amount\":"
		out.RawString(prefix)
		if in.TotalAmount == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.TotalAmount {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalTinyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ReleaseStats) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson16010b1EncodeGithubComCosmWasmCosmwasmGoExampleVoterSrcState(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v ReleaseStats) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson16010b1EncodeGithubComCosmWasmCosmwasmGoExampleVoterSrcState(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ReleaseStats) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson16010b1DecodeGithubComCosmWasmCosmwasmGoExampleVoterSrcState(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *ReleaseStats) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson16010b1DecodeGithubComCosmWasmCosmwasmGoExampleVoterSrcState(l, v)
}
