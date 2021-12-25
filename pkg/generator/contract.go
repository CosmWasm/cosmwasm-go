package generator

import (
	"encoding/json"
	"fmt"
	"github.com/CosmWasm/tinyjson"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	QueryHandlerPrefix     = "Query"
	ExecHandlerPrefix      = "Exec"
	MigrateHandlerName     = "Migrate"
	InstantiateHandlerName = "Instantiate"
)

const (
	typesPkg = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/std/types")
)

var (
	zstStruct           = reflect.TypeOf(struct{}{})
	depsType            = reflect.TypeOf((*std.Deps)(nil))
	envType             = reflect.TypeOf((*types.Env)(nil))
	infoType            = reflect.TypeOf((*types.MessageInfo)(nil))
	tinyjsonMarshaler   = reflect.TypeOf((*tinyjson.Marshaler)(nil)).Elem()
	tinyjsonUnmarshaler = reflect.TypeOf((*tinyjson.Unmarshaler)(nil)).Elem()
	jsonMarshaler       = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	jsonUnmarshaler     = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	errType             = reflect.TypeOf((*error)(nil)).Elem()
)

type MigrateDescriptor struct {
	Foreign    bool
	InputType  protogen.GoIdent
	MethodName string
}

type InstantiateDescriptor struct {
	Foreign    bool
	InputType  protogen.GoIdent
	MethodName string
}

type ExecDescriptor struct {
	JSONName    string
	InUnionName string
	InputType   protogen.GoIdent
	Foreign     bool
	MethodName  string
}

type QueryDescriptor struct {
	JSONName    string
	InUnionName string
	InputType   protogen.GoIdent
	OutputType  protogen.GoIdent
	Foreign     bool // imported or not
	MethodName  string
}

func NewContract(pkg string, v interface{}) (*Contract, error) {
	typ := reflect.TypeOf(v)
	if !isStateless(typ) {
		return nil, fmt.Errorf("contract must be a stateless structure which means it can only be struct{} or have fields or embed with struct{}")
	}

	return &Contract{
		Generator:   NewGenerator(),
		typ:         typ,
		exec:        map[string]ExecDescriptor{},
		query:       map[string]QueryDescriptor{},
		instantiate: nil,
		pkg:         pkg,
		tinyjsonGen: nil,
	}, nil
}

// Contract takes care of generating the contract boilerplate code.
type Contract struct {
	*Generator
	typ         reflect.Type
	exec        map[string]ExecDescriptor
	query       map[string]QueryDescriptor
	instantiate *InstantiateDescriptor
	migrate     *MigrateDescriptor

	pkg         string
	tinyjsonGen []string // types that need to have a tinyjson impl
}

func (g *Contract) Generate() error {
	err := g.process()
	if err != nil {
		return err
	}

	err = g.genPkg()
	if err != nil {
		return err
	}

	err = g.genExecute()
	if err != nil {
		return err
	}

	err = g.genQuery()
	if err != nil {
		return err
	}

	err = g.genInstantiate()
	if err != nil {
		return err
	}

	err = g.genMigrate()
	if err != nil {
		return err
	}
	return nil
}

func (g *Contract) WriteTo(path string) (err error) {
	// write contract boilerplate
	err = g.Generator.WriteTo(path)
	if err != nil {
		return err
	}
	// in case following steps fail, then we remove the generated file
	defer func() {
		if err != nil {
			_ = os.Remove(path)
		}
	}()

	// write tinyjson placeholders
	tinyjsonFilePath := filepath.Join(
		filepath.Dir(path),
		strings.TrimSuffix(filepath.Base(path), ".go")+"_tinyjson.go",
	)
	f, err := os.OpenFile(tinyjsonFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	tinyjsonGen := NewGenerator()
	tinyjsonGen.P("package ", g.pkg)
	const (
		tinyjsonPkg = protogen.GoImportPath("github.com/CosmWasm/tinyjson")
		jlexerPkg   = protogen.GoImportPath("github.com/CosmWasm/tinyjson/jlexer")
		jwriterPkg  = protogen.GoImportPath("github.com/CosmWasm/tinyjson/jwriter")
	)

	tinyjsonGen.P("var (")
	tinyjsonGen.P("_ *", jlexerPkg.Ident("Lexer"))
	tinyjsonGen.P("_ *", jwriterPkg.Ident("Writer"))
	tinyjsonGen.P("_ ", tinyjsonPkg.Ident("Marshaler"))
	tinyjsonGen.P(")")
	for _, typeName := range g.tinyjsonGen {
		tinyjsonGen.P("func (x *", typeName, ") MarshalJSON() ([]byte,error) { panic(0) }")
		tinyjsonGen.P("func (x *", typeName, ") MarshalTinyJSON(_ *", jwriterPkg.Ident("Writer"), ") { panic(0) }")
		tinyjsonGen.P("func (x *", typeName, ") UnmarshalJSON(b []byte) error { panic(0) }")
		tinyjsonGen.P("func (x *", typeName, ") UnmarshalTinyJSON(_ *", jlexerPkg.Ident("Lexer"), ") { panic(0) }")
	}

	content, err := tinyjsonGen.Content()
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func isReply(name string) bool {
	// todo
	return false
}

func isMigrate(name string) bool {
	return name == MigrateHandlerName
}

func isSudo(name string) bool {
	// TODO
	return false
}

func (g *Contract) addQuery(method reflect.Method) error {
	funcType := method.Func.Type()

	inputs := make([]reflect.Type, funcType.NumIn()-1)
	for i := 1; i < funcType.NumIn(); i++ {
		inputs[i-1] = funcType.In(i)
	}

	outputs := make([]reflect.Type, funcType.NumOut())
	for i := 0; i < funcType.NumOut(); i++ {
		outputs[i] = funcType.Out(i)
	}

	// assert query inputs
	if len(inputs) != 3 {
		return fmt.Errorf("invalid query method %s inputs, expected 3 (std.Deps, types.Env, queryMsg), got: %d", method.Name, len(inputs))
	}

	if depsType != inputs[0] {
		return fmt.Errorf("first input must be %s, got %s", depsType, inputs[0])
	}

	if envType != inputs[1] {
		return fmt.Errorf("second input must be %s, got %s", envType, inputs[1])
	}

	if inputs[2].Kind() != reflect.Ptr {
		return fmt.Errorf("third input must be a pointer, got: %s", inputs[2].Kind())
	}

	if inputs[2].Elem().Kind() != reflect.Struct {
		return fmt.Errorf("third input must be a struct pointer, got: %s", inputs[2].Elem().Kind())
	}

	// we need to make sure it implements json unmarshaler
	if !implementsUnmarshaler(inputs[2]) {
		return fmt.Errorf("third input must implement tinyjson.Unmarshaler or json.Unmarshaler")
	}

	inputPkg := protogen.GoImportPath(inputs[2].Elem().PkgPath())

	// assert outputs
	if len(outputs) != 2 {
		return fmt.Errorf("unexpected number of query outputs, want 2 (queryMsgResponse, error); got %d", len(outputs))
	}

	if outputs[0].Kind() != reflect.Ptr {
		return fmt.Errorf("first output must be a pointer, got: %s", outputs[0].Kind())
	}
	if outputs[0].Elem().Kind() != reflect.Struct {
		return fmt.Errorf("first output must be a struct pointer: %s", outputs[0].Elem().Kind())
	}

	if !implementsMarshaler(outputs[0]) {
		return fmt.Errorf("first output must implement json marshaler")
	}

	if outputs[1] != errType {
		return fmt.Errorf("second output must be error, got: %s", outputs[1])
	}

	outputPkg := protogen.GoImportPath(outputs[0].Elem().PkgPath())

	inUnionName := strings.TrimPrefix(method.Name, QueryHandlerPrefix)
	jsonName := strcase.ToSnake(inUnionName)

	g.query[jsonName] = QueryDescriptor{
		JSONName:    jsonName,
		InUnionName: inUnionName,
		InputType:   inputPkg.Ident(inputs[2].Elem().Name()),
		OutputType:  outputPkg.Ident(outputs[0].Elem().Name()),
		MethodName:  method.Name,
		Foreign:     g.typ.PkgPath() != inputs[2].Elem().PkgPath(),
	}

	return nil
}

func (g *Contract) addExec(method reflect.Method) error {
	funcType := method.Func.Type()
	inputs := make([]reflect.Type, funcType.NumIn()-1)
	outputs := make([]reflect.Type, funcType.NumOut())

	for i := 1; i < funcType.NumIn(); i++ {
		inputs[i-1] = funcType.In(i)
	}

	for i := 0; i < funcType.NumOut(); i++ {
		outputs[i] = funcType.Out(i)
	}

	if len(inputs) != 4 {
		return fmt.Errorf("unexpected number of inputs, expected 4: std.Deps, types.Env, types.MessageInfo, message")
	}

	// assert types
	if inputs[0] != depsType {
		return fmt.Errorf("first input must be of type *std.Deps")
	}
	if inputs[1] != envType {
		return fmt.Errorf("second input must be of type *types.Env")
	}
	if inputs[2] != infoType {
		return fmt.Errorf("third input must be of type *types.MessageInfo")
	}

	if inputs[3].Kind() != reflect.Ptr {
		return fmt.Errorf("fourth input must be a pointer")
	}

	if !implementsUnmarshaler(inputs[3]) {
		return fmt.Errorf("fourth input must implement tinyjson.Unmarshaler")
	}

	inputPkg := protogen.GoImportPath(inputs[3].PkgPath())

	inUnionName := strings.TrimPrefix(method.Name, ExecHandlerPrefix)
	jsonName := strcase.ToSnake(inUnionName)
	g.exec[jsonName] = ExecDescriptor{
		JSONName:    jsonName,
		InUnionName: inUnionName,
		InputType:   inputPkg.Ident(inputs[3].Elem().Name()),
		Foreign:     inputs[3].Elem().PkgPath() != g.typ.PkgPath(),
		MethodName:  method.Name,
	}

	return nil
}

func (g *Contract) addSudo(method reflect.Method) error {
	// todo
	return nil
}

func (g *Contract) addMigrate(method reflect.Method) error {
	funcType := method.Func.Type()
	inputs := make([]reflect.Type, funcType.NumIn()-1)
	outputs := make([]reflect.Type, funcType.NumOut())

	for i := 1; i < funcType.NumIn(); i++ {
		inputs[i-1] = funcType.In(i)
	}

	for i := 0; i < funcType.NumOut(); i++ {
		outputs[i] = funcType.Out(i)
	}

	if len(inputs) != 3 {
		return fmt.Errorf("unexpected number of inputs, expected 3: std.Deps, types.Env, message")
	}

	// assert types
	if inputs[0] != depsType {
		return fmt.Errorf("first input must be of type *std.Deps")
	}
	if inputs[1] != envType {
		return fmt.Errorf("second input must be of type *types.Env")
	}

	if inputs[2].Kind() != reflect.Ptr {
		return fmt.Errorf("third input must be a pointer")
	}

	if !implementsUnmarshaler(inputs[2]) {
		return fmt.Errorf("third input must implement tinyjson.Unmarshaler")
	}

	inputPkg := protogen.GoImportPath(inputs[2].Elem().PkgPath())

	g.migrate = &MigrateDescriptor{
		InputType:  inputPkg.Ident(inputs[2].Elem().Name()),
		Foreign:    inputs[2].Elem().PkgPath() != g.typ.PkgPath(),
		MethodName: method.Name,
	}

	return nil
}

func (g *Contract) addReply(method reflect.Method) error {
	// todo
	return nil
}

func (g *Contract) process() error {
	// extract contract methods from the type
	for i := 0; i < g.typ.NumMethod(); i++ {
		method := g.typ.Method(i)
		// we ignore unexported methods.
		if !method.IsExported() {
			continue
		}

		switch {
		case isQuery(method.Name):
			err := g.addQuery(method)
			if err != nil {
				return err
			}
		case isExec(method.Name):
			err := g.addExec(method)
			if err != nil {
				return err
			}
		case isSudo(method.Name):
			err := g.addSudo(method)
			if err != nil {
				return err
			}
		case isMigrate(method.Name):
			err := g.addMigrate(method)
			if err != nil {
				return err
			}
		case isReply(method.Name):
			err := g.addReply(method)
			if err != nil {
				return err
			}
		case isInstantiate(method.Name):
			err := g.addInstantiate(method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isInstantiate(name string) bool {
	return name == InstantiateHandlerName
}

func (g *Contract) genQuery() error {
	err := g.genQueryUnion()
	if err != nil {
		return err
	}

	err = g.genQueryHandler()
	if err != nil {
		return err
	}

	err = g.genQueryMsgHelpers()
	if err != nil {

	}
	return nil
}

func (g *Contract) genQueryUnion() error {
	g.P("// QueryMsg is the union type used to process queries towards the contract.")
	g.P("type QueryMsg struct {")
	for jsonName, desc := range g.query {
		switch desc.Foreign {
		case true:
			g.P(desc.InUnionName, " ", "*", desc.InputType, " `json:\"", jsonName, "\"`")
		case false:
			g.P(desc.InUnionName, " ", "*", desc.InputType.GoName, " `json:\"", jsonName, "\"`")
		}
	}
	g.P("}")
	g.P()

	g.addTinyJSONImpl("QueryMsg")

	return nil
}

func (g *Contract) genQueryHandler() error {
	g.P("func Query(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", queryBytes []byte) ([]byte, error) {")
	g.P("query := new(QueryMsg)")
	g.P("err := query.UnmarshalJSON(queryBytes)")
	g.P("if err != nil { return nil, err }")
	g.P("switch {")
	for _, desc := range g.query {
		g.P("case query.", desc.InUnionName, " != nil:")
		g.P("resp, err := ", g.typ.Name(), "{}.", desc.MethodName, "(deps, &env, query.", desc.InUnionName, ")")
		g.P("if err != nil { return nil, err }")
		g.P("return resp.MarshalJSON()")
	}
	g.P("default:")
	g.P("panic(1)") // TODO we need an error pkg for common errors...
	g.P("}")
	g.P("}")
	g.P()
	return nil
}

func (g *Contract) genPkg() error {
	g.P("package ", g.pkg)
	g.P()

	return nil
}

func (g *Contract) genQueryMsgHelpers() error {
	for _, desc := range g.query {
		if desc.Foreign {
			continue
		}

		g.P("func (x *", desc.InputType.GoName, ") AsQueryMsg() *QueryMsg {")
		g.P("return &QueryMsg{", desc.InUnionName, ": x}")
		g.P("}")
		g.P()
	}

	return nil
}

func (g *Contract) addTinyJSONImpl(s string) {
	g.tinyjsonGen = append(g.tinyjsonGen, s)
}

func (g *Contract) genExecute() error {
	err := g.genExecuteUnion()
	if err != nil {
		return err
	}

	err = g.genExecuteHandler()
	if err != nil {
		return err
	}

	err = g.genExecuteHelpers()
	if err != nil {
		return err
	}

	return nil
}

func (g *Contract) genExecuteUnion() error {
	g.P("// ExecuteMsg is the union type used to process execution messages towards the contract.")
	g.P("type ExecuteMsg struct {")
	for jsonName, desc := range g.exec {
		switch desc.Foreign {
		case true:
			g.P(desc.InUnionName, " ", "*", desc.InputType, " `json:\"", jsonName, "\"`")
		case false:
			g.P(desc.InUnionName, " ", "*", desc.InputType.GoName, " `json:\"", jsonName, "\"`")
		}
	}
	g.P("}")
	g.P()

	g.addTinyJSONImpl("ExecuteMsg")

	return nil
}

func (g *Contract) genExecuteHandler() error {
	g.P("func Execute(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", info ", typesPkg.Ident("MessageInfo"), ", messageBytes []byte) (*", typesPkg.Ident("Response"), ", error) {")
	g.P("msg := new(ExecuteMsg)")
	g.P("err := msg.UnmarshalJSON(messageBytes)")
	g.P("if err != nil { return nil, err }")
	g.P("switch {")
	for _, desc := range g.exec {
		g.P("case msg.", desc.InUnionName, " != nil:")
		g.P("resp, err := ", g.typ.Name(), "{}.", desc.MethodName, "(deps, &env, &info, msg.", desc.InUnionName, ")")
		g.P("if err != nil { return nil, err }")
		g.P("return resp, nil")
	}
	g.P("default:")
	g.P("panic(1)") // TODO we need an error pkg for common errors...
	g.P("}")
	g.P("}")
	g.P()
	return nil
}

func (g *Contract) genExecuteHelpers() error {
	return nil
}

func (g *Contract) genInstantiate() error {
	g.P("func Instantiate(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", info ", typesPkg.Ident("MessageInfo"), ", instantiateBytes []byte)", "(*", typesPkg.Ident("Response"), ",error) {")
	if g.instantiate == nil {
		g.P("return &", typesPkg.Ident("Response"), "{}, nil")
		g.P("}")
		return nil
	}
	switch g.instantiate.Foreign {
	case true:
		g.P("initMsg := new(", g.instantiate.InputType, ")")
	case false:
		g.P("initMsg := new(", g.instantiate.InputType.GoName, ")")
	}
	g.P("err := initMsg.UnmarshalJSON(instantiateBytes)")
	g.P("if err != nil { return nil, err }")
	g.P("return ", g.typ.Name(), "{}.", g.instantiate.MethodName, "(deps, &env, &info, initMsg)")
	g.P("}")
	g.P()
	return nil
}

func (g *Contract) addInstantiate(method reflect.Method) error {
	funcType := method.Func.Type()
	inputs := make([]reflect.Type, funcType.NumIn()-1)
	outputs := make([]reflect.Type, funcType.NumOut())

	for i := 1; i < funcType.NumIn(); i++ {
		inputs[i-1] = funcType.In(i)
	}

	for i := 0; i < funcType.NumOut(); i++ {
		outputs[i] = funcType.Out(i)
	}

	if len(inputs) != 4 {
		return fmt.Errorf("unexpected number of inputs, expected 4: std.Deps, types.Env, types.MessageInfo, message")
	}

	// assert types
	if inputs[0] != depsType {
		return fmt.Errorf("first input must be of type *std.Deps")
	}
	if inputs[1] != envType {
		return fmt.Errorf("second input must be of type *types.Env")
	}
	if inputs[2] != infoType {
		return fmt.Errorf("third input must be of type *types.MessageInfo")
	}

	if inputs[3].Kind() != reflect.Ptr {
		return fmt.Errorf("fourth input must be a pointer")
	}

	if !implementsUnmarshaler(inputs[3]) {
		return fmt.Errorf("fourth input must implement tinyjson.Unmarshaler")
	}

	inputPkg := protogen.GoImportPath(inputs[3].Elem().PkgPath())

	g.instantiate = &InstantiateDescriptor{
		Foreign:    inputs[3].Elem().PkgPath() != g.typ.PkgPath(),
		InputType:  inputPkg.Ident(inputs[3].Elem().Name()),
		MethodName: method.Name,
	}

	return nil
}

func (g *Contract) genMigrate() error {
	g.P("func Migrate(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", migrateBytes []byte) (*", typesPkg.Ident("Response"), ", error) {")
	if g.migrate == nil {
		g.P("return &", typesPkg.Ident("Response"), "{}, nil")
		g.P("}")
		return nil
	}

	switch g.migrate.Foreign {
	case true:
		g.P("msg := new(", g.migrate.InputType, ")")
	case false:
		g.P("msg := new(", g.migrate.InputType.GoName, ")")
	}
	g.P("err := msg.UnmarshalJSON(migrateBytes)")
	g.P("if err != nil { return nil, err }")
	g.P("return ", g.typ.Name(), "{}.", g.migrate.MethodName, "(deps, &env, msg)")
	g.P("}")
	return nil
}

func isQuery(name string) bool {
	return strings.HasPrefix(name, QueryHandlerPrefix)
}

func isExec(name string) bool {
	return strings.HasPrefix(name, ExecHandlerPrefix)
}

func implementsUnmarshaler(t reflect.Type) bool {
	return t.Implements(jsonUnmarshaler) || t.Implements(tinyjsonUnmarshaler)
}

func implementsMarshaler(t reflect.Type) bool {
	return t.Implements(jsonMarshaler) || t.Implements(tinyjsonMarshaler)
}

// isStateless checks if the provided type is stateless
// which means all the fields (if any) must be zero sized struct, or contain only
// zero sized structs fields.
func isStateless(t reflect.Type) bool {
	if t.ConvertibleTo(zstStruct) {
		return true
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !isStateless(field.Type) {
			return false
		}
	}

	return true
}
