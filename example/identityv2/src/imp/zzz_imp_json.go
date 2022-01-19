// Code generated by tinyjson for marshaling/unmarshaling. DO NOT EDIT.

package imp

import (
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

func tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp(in *jlexer.Lexer, out *ImportedQueryResponse) {
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
func tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp(out *jwriter.Writer, in ImportedQueryResponse) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ImportedQueryResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v ImportedQueryResponse) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ImportedQueryResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *ImportedQueryResponse) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp(l, v)
}
func tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp1(in *jlexer.Lexer, out *ImportedQuery) {
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
func tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp1(out *jwriter.Writer, in ImportedQuery) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ImportedQuery) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v ImportedQuery) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ImportedQuery) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp1(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *ImportedQuery) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp1(l, v)
}
func tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp2(in *jlexer.Lexer, out *ImportedMessage) {
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
func tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp2(out *jwriter.Writer, in ImportedMessage) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ImportedMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v ImportedMessage) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonFb78fff9EncodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ImportedMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp2(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *ImportedMessage) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonFb78fff9DecodeGithubComCosmwasmCosmwasmGoExampleIdentityv2SrcImp2(l, v)
}