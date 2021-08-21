// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package std

import (
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd(in *jlexer.Lexer, out *MessageInfo) {
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
		case "sender":
			out.Sender = string(in.String())
		case "funds":
			if in.IsNull() {
				in.Skip()
				out.Funds = nil
			} else {
				in.Delim('[')
				if out.Funds == nil {
					if !in.IsDelim(']') {
						out.Funds = make([]Coin, 0, 2)
					} else {
						out.Funds = []Coin{}
					}
				} else {
					out.Funds = (out.Funds)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Coin
					(v1).UnmarshalEasyJSON(in)
					out.Funds = append(out.Funds, v1)
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
func easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd(out *jwriter.Writer, in MessageInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"sender\":"
		out.RawString(prefix[1:])
		out.String(string(in.Sender))
	}
	{
		const prefix string = ",\"funds\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v2, v3 := range in.Funds {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd(l, v)
}
func easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd1(in *jlexer.Lexer, out *Env) {
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
		case "block":
			(out.Block).UnmarshalEasyJSON(in)
		case "contract":
			(out.Contract).UnmarshalEasyJSON(in)
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
func easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd1(out *jwriter.Writer, in Env) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"block\":"
		out.RawString(prefix[1:])
		(in.Block).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"contract\":"
		out.RawString(prefix)
		(in.Contract).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Env) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Env) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Env) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Env) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd1(l, v)
}
func easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd2(in *jlexer.Lexer, out *ContractInfo) {
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
		case "address":
			out.Address = string(in.String())
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
func easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd2(out *jwriter.Writer, in ContractInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix[1:])
		out.String(string(in.Address))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ContractInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ContractInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ContractInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ContractInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd2(l, v)
}
func easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd3(in *jlexer.Lexer, out *BlockInfo) {
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
		case "height":
			out.Height = uint64(in.Uint64())
		case "time":
			out.Time = uint64(in.Uint64Str())
		case "chain_id":
			out.ChainID = string(in.String())
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
func easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd3(out *jwriter.Writer, in BlockInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"height\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Height))
	}
	{
		const prefix string = ",\"time\":"
		out.RawString(prefix)
		out.Uint64Str(uint64(in.Time))
	}
	{
		const prefix string = ",\"chain_id\":"
		out.RawString(prefix)
		out.String(string(in.ChainID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BlockInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BlockInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3ba44563EncodeGithubComCosmwasmCosmwasmGoStd3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BlockInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BlockInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3ba44563DecodeGithubComCosmwasmCosmwasmGoStd3(l, v)
}
