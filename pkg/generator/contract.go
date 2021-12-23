package generator

import (
	"encoding/json"
	"fmt"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	"io"
	"reflect"
	"strings"
)

const (
	QueryHandlerPrefix = "Query"
	ExecHandlerPrefix  = "Exec"
)

const (
	typesPkg = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/std/types")
)

var (
	zstStruct          = reflect.TypeOf(struct{}{})
	depsType           = reflect.TypeOf((*std.Deps)(nil))
	envType            = reflect.TypeOf((*types.Env)(nil))
	infoType           = reflect.TypeOf((*types.MessageInfo)(nil))
	marshalInterface   = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	unmarshalInterface = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	errType            = reflect.TypeOf((*error)(nil)).Elem()
)

type ExecDescriptor struct {
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
	if !typ.ConvertibleTo(zstStruct) {
		return nil, fmt.Errorf("%s is not convertible to struct{}, contracts must be stateless", typ.Name())
	}

	return &Contract{
		Generator: NewGenerator(),
		typ:       typ,
		exec:      map[string]ExecDescriptor{},
		query:     map[string]QueryDescriptor{},
		pkg:       pkg,
	}, nil
}

// Contract takes care of generating the contract boilerplate code.
type Contract struct {
	*Generator
	typ   reflect.Type
	exec  map[string]ExecDescriptor
	query map[string]QueryDescriptor
	pkg   string
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

	err = g.genQuery()
	if err != nil {
		return err
	}

	return nil
}

func (g *Contract) Write(f io.Writer) error {
	content, err := g.Content()
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	return err
}

func isReply(name string) bool {
	// todo
	return false
}

func isMigrate(name string) bool {
	// todo
	return false
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
	if !inputs[2].Implements(unmarshalInterface) {
		return fmt.Errorf("third input must implement json.Unmarshaler")
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

	if !outputs[0].Implements(marshalInterface) {
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
	return nil
}

func (g *Contract) addSudo(method reflect.Method) error {
	// todo
	return nil
}

func (g *Contract) addMigrate(method reflect.Method) error {
	// todo
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
		}
	}
	return nil
}

func (g *Contract) genQuery() error {
	err := g.genQueryUnion()
	if err != nil {
		return err
	}

	err = g.genQueryExec()
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

	// TODO implement json unmarshal
	g.P("func (x *QueryMsg) UnmarshalJSON(b []byte) error {")
	g.P("panic(0)")
	g.P("}")

	// TODO implement json marshal
	g.P("func (x *QueryMsg) MarshalJSON() ([]byte, error) {")
	g.P("panic(0)")
	g.P("}")

	g.P()
	return nil
}

func (g *Contract) genQueryExec() error {
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
	}

	return nil
}

func isQuery(name string) bool {
	return strings.HasPrefix(name, QueryHandlerPrefix)
}

func isExec(name string) bool {
	return strings.HasPrefix(name, ExecHandlerPrefix)
}
