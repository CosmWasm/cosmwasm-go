package src

import (
	tinyjson "github.com/CosmWasm/tinyjson"
	jlexer "github.com/CosmWasm/tinyjson/jlexer"
	jwriter "github.com/CosmWasm/tinyjson/jwriter"
)

var (
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ tinyjson.Marshaler
)

func tinyjsonDd15385dDecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(in *jlexer.Lexer, out *QueryMsg) {
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
		case "identity":
			if in.IsNull() {
				in.Skip()
				out.Identity = nil
			} else {
				if out.Identity == nil {
					out.Identity = new(QueryIdentity)
				}
				(*out.Identity).UnmarshalTinyJSON(in)
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
func tinyjsonDd15385dEncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(out *jwriter.Writer, in QueryMsg) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"identity\":"
		out.RawString(prefix[1:])
		if in.Identity == nil {
			out.RawString("null")
		} else {
			(*in.Identity).MarshalTinyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueryMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonDd15385dEncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v QueryMsg) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonDd15385dEncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueryMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonDd15385dDecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *QueryMsg) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonDd15385dDecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(l, v)
}
func tinyjsonDd15385dDecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(in *jlexer.Lexer, out *ExecuteMsg) {
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
		case "create_identity":
			if in.IsNull() {
				in.Skip()
				out.CreateIdentity = nil
			} else {
				if out.CreateIdentity == nil {
					out.CreateIdentity = new(MsgCreateIdentity)
				}
				(*out.CreateIdentity).UnmarshalTinyJSON(in)
			}
		case "delete_identity":
			if in.IsNull() {
				in.Skip()
				out.DeleteIdentity = nil
			} else {
				if out.DeleteIdentity == nil {
					out.DeleteIdentity = new(MsgDelete)
				}
				(*out.DeleteIdentity).UnmarshalTinyJSON(in)
			}
		case "update_city":
			if in.IsNull() {
				in.Skip()
				out.UpdateCity = nil
			} else {
				if out.UpdateCity == nil {
					out.UpdateCity = new(MsgUpdateCity)
				}
				(*out.UpdateCity).UnmarshalTinyJSON(in)
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
func tinyjsonDd15385dEncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(out *jwriter.Writer, in ExecuteMsg) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"create_identity\":"
		out.RawString(prefix[1:])
		if in.CreateIdentity == nil {
			out.RawString("null")
		} else {
			(*in.CreateIdentity).MarshalTinyJSON(out)
		}
	}
	{
		const prefix string = ",\"delete_identity\":"
		out.RawString(prefix)
		if in.DeleteIdentity == nil {
			out.RawString("null")
		} else {
			(*in.DeleteIdentity).MarshalTinyJSON(out)
		}
	}
	{
		const prefix string = ",\"update_city\":"
		out.RawString(prefix)
		if in.UpdateCity == nil {
			out.RawString("null")
		} else {
			(*in.UpdateCity).MarshalTinyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExecuteMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonDd15385dEncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v ExecuteMsg) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonDd15385dEncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExecuteMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonDd15385dDecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *ExecuteMsg) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonDd15385dDecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(l, v)
}
