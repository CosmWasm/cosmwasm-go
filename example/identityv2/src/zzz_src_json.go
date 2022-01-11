// Code generated by tinyjson for marshaling/unmarshaling. DO NOT EDIT.

package src

import (
	tinyjson "github.com/CosmWasm/tinyjson"
	jlexer "github.com/CosmWasm/tinyjson/jlexer"
	jwriter "github.com/CosmWasm/tinyjson/jwriter"
	imp "github.com/cosmwasm/cosmwasm-go/example/identityv2/src/imp"
)

// suppress unused package warning
var (
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ tinyjson.Marshaler
)

func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(in *jlexer.Lexer, out *QueryMsg) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "query_identity":
			if in.IsNull() {
				in.Skip()
				out.QueryIdentity = nil
			} else {
				if out.QueryIdentity == nil {
					out.QueryIdentity = new(QueryIdentity)
				}
				(*out.QueryIdentity).UnmarshalTinyJSON(in)
			}
		case "query_imported":
			if in.IsNull() {
				in.Skip()
				out.QueryImported = nil
			} else {
				if out.QueryImported == nil {
					out.QueryImported = new(imp.ImportedQuery)
				}
				(*out.QueryImported).UnmarshalTinyJSON(in)
			}
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(out *jwriter.Writer, in QueryMsg) {
	out.RawByte('{')
	first := true
	_ = first
	if in.QueryIdentity != nil {
		const prefix string = ",\"query_identity\":"
		first = false
		out.RawString(prefix[1:])
		(*in.QueryIdentity).MarshalTinyJSON(out)
	}
	if in.QueryImported != nil {
		const prefix string = ",\"query_imported\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.QueryImported).MarshalTinyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueryMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v QueryMsg) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueryMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *QueryMsg) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(in *jlexer.Lexer, out *QueryIdentityResponse) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "person":
			if in.IsNull() {
				in.Skip()
				out.Person = nil
			} else {
				if out.Person == nil {
					out.Person = new(Person)
				}
				(*out.Person).UnmarshalTinyJSON(in)
			}
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(out *jwriter.Writer, in QueryIdentityResponse) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Person != nil {
		const prefix string = ",\"person\":"
		first = false
		out.RawString(prefix[1:])
		(*in.Person).MarshalTinyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueryIdentityResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v QueryIdentityResponse) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueryIdentityResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *QueryIdentityResponse) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src1(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src2(in *jlexer.Lexer, out *QueryIdentity) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src2(out *jwriter.Writer, in QueryIdentity) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueryIdentity) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v QueryIdentity) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueryIdentity) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src2(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *QueryIdentity) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src2(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src3(in *jlexer.Lexer, out *Person) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "address":
			out.Address = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "city":
			out.City = string(in.String())
		case "postal_code":
			out.PostalCode = int32(in.Int32())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src3(out *jwriter.Writer, in Person) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Address != "" {
		const prefix string = ",\"address\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Address))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Surname != "" {
		const prefix string = ",\"surname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Surname))
	}
	if in.City != "" {
		const prefix string = ",\"city\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.City))
	}
	if in.PostalCode != 0 {
		const prefix string = ",\"postal_code\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.PostalCode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Person) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v Person) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Person) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src3(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *Person) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src3(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src4(in *jlexer.Lexer, out *MsgUpdateCity) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "city":
			out.City = string(in.String())
		case "postal_code":
			out.PostalCode = int32(in.Int32())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src4(out *jwriter.Writer, in MsgUpdateCity) {
	out.RawByte('{')
	first := true
	_ = first
	if in.City != "" {
		const prefix string = ",\"city\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.City))
	}
	if in.PostalCode != 0 {
		const prefix string = ",\"postal_code\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.PostalCode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MsgUpdateCity) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v MsgUpdateCity) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MsgUpdateCity) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src4(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *MsgUpdateCity) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src4(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src5(in *jlexer.Lexer, out *MsgMigrate) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src5(out *jwriter.Writer, in MsgMigrate) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MsgMigrate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v MsgMigrate) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MsgMigrate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src5(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *MsgMigrate) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src5(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src6(in *jlexer.Lexer, out *MsgInstantiate) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src6(out *jwriter.Writer, in MsgInstantiate) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MsgInstantiate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v MsgInstantiate) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MsgInstantiate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src6(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *MsgInstantiate) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src6(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src7(in *jlexer.Lexer, out *MsgDelete) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src7(out *jwriter.Writer, in MsgDelete) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MsgDelete) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v MsgDelete) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MsgDelete) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src7(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *MsgDelete) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src7(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src8(in *jlexer.Lexer, out *MsgCreateIdentity) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "city":
			out.City = string(in.String())
		case "postal_code":
			out.PostalCode = int32(in.Int32())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src8(out *jwriter.Writer, in MsgCreateIdentity) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.Surname != "" {
		const prefix string = ",\"surname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Surname))
	}
	if in.City != "" {
		const prefix string = ",\"city\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.City))
	}
	if in.PostalCode != 0 {
		const prefix string = ",\"postal_code\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.PostalCode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MsgCreateIdentity) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v MsgCreateIdentity) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MsgCreateIdentity) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src8(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *MsgCreateIdentity) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src8(l, v)
}
func tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src9(in *jlexer.Lexer, out *ExecuteMsg) {
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
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "imported_message":
			if in.IsNull() {
				in.Skip()
				out.ImportedMessage = nil
			} else {
				if out.ImportedMessage == nil {
					out.ImportedMessage = new(imp.ImportedMessage)
				}
				(*out.ImportedMessage).UnmarshalTinyJSON(in)
			}
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
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src9(out *jwriter.Writer, in ExecuteMsg) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ImportedMessage != nil {
		const prefix string = ",\"imported_message\":"
		first = false
		out.RawString(prefix[1:])
		(*in.ImportedMessage).MarshalTinyJSON(out)
	}
	if in.CreateIdentity != nil {
		const prefix string = ",\"create_identity\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.CreateIdentity).MarshalTinyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExecuteMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v ExecuteMsg) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjson97bc4d59EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExecuteMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src9(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *ExecuteMsg) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjson97bc4d59DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2Src9(l, v)
}
