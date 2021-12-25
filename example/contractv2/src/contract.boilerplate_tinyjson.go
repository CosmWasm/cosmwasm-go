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

func (x *ExecuteMsg) MarshalJSON() ([]byte, error)      { panic(0) }
func (x *ExecuteMsg) MarshalTinyJSON(_ *jwriter.Writer) { panic(0) }
func (x *ExecuteMsg) UnmarshalJSON(b []byte) error      { panic(0) }
func (x *ExecuteMsg) UnmarshalTinyJSON(_ *jlexer.Lexer) { panic(0) }
func (x *QueryMsg) MarshalJSON() ([]byte, error)        { panic(0) }
func (x *QueryMsg) MarshalTinyJSON(_ *jwriter.Writer)   { panic(0) }
func (x *QueryMsg) UnmarshalJSON(b []byte) error        { panic(0) }
func (x *QueryMsg) UnmarshalTinyJSON(_ *jlexer.Lexer)   { panic(0) }
