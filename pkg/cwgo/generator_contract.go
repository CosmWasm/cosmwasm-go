package cwgo

import (
	"github.com/cosmwasm/cosmwasm-go/pkg/cwgo/gen"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	typesPkg = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/std/types")
)

type contractGen struct {
	*gen.Gen

	dir, pkg string

	typeName    string
	execs       map[string]*Exec
	queries     map[string]*Query
	migrate     *Migrate
	instantiate *Instantiate
}

func newContractGen(dir, pkg string, contract *Contract) *contractGen {
	return &contractGen{
		Gen:         gen.NewGen(dir, contractFileName(pkg, contract.Name)),
		dir:         dir,
		pkg:         pkg,
		typeName:    contract.Name,
		execs:       contract.Execs,
		queries:     contract.Queries,
		migrate:     contract.Migrate,
		instantiate: contract.Instantiate,
	}
}

func (g *contractGen) generate() error {

	// gen pkg
	g.P("package ", g.pkg)

	// gen union
	g.unions()
	// gen type helpers
	g.typeHelpers()
	// gen Execute function
	g.execFunc()
	// gen Query function
	g.queryFunc()
	// gen Migrate function
	g.migrateFunc()
	// gen Instantiate function
	g.instantiateFunc()

	return g.Generate()
}

func (g *contractGen) unions() {
	g.P("type QueryMsg struct {")
	for name, query := range g.queries {
		var typ interface{}
		switch query.Input.ImportPath {
		case "":
			typ = query.Input.Name
		default:
			typ = protogen.GoIdent{
				GoName:       query.Input.Name,
				GoImportPath: query.Input.ImportPath,
			}
		}
		g.P(name, " *", typ, "`json:\"", strcase.ToSnake(name), "\"`")
	}
	g.P("}")
	g.P()
	g.P("type ExecuteMsg struct {")
	for name, exec := range g.execs {
		var typ interface{}
		switch exec.Input.ImportPath {
		case "":
			typ = exec.Input.Name
		default:
			typ = protogen.GoIdent{
				GoName:       exec.Input.Name,
				GoImportPath: exec.Input.ImportPath,
			}
		}
		g.P(name, " *", typ, "`json:\"", strcase.ToSnake(name), "\"`")
	}
	g.P("}")
	g.P()
}

func (g *contractGen) typeHelpers() {
	for _, query := range g.queries {
		// skip imported types
		if query.Input.ImportPath != "" {
			continue
		}
		g.P("func (x *", query.Input.Name, ") AsQueryMsg() QueryMsg {")
		g.P("return QueryMsg{", query.MethodName, ": x,}")
		g.P("}")
		g.P()
	}

	for _, exec := range g.execs {
		if exec.Input.ImportPath != "" {
			continue
		}

		g.P("func (x *", exec.Input.Name, ") AsExecuteMsg() ExecuteMsg {")
		g.P("return ExecuteMsg{", exec.MethodName, ": x,}")
		g.P("}")
		g.P()
	}
}

func (g *contractGen) execFunc() {
	g.P("func Execute(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", info ", typesPkg.Ident("MessageInfo"), ", messageBytes []byte) (*", typesPkg.Ident("Response"), ", error) {")
	g.P("msg := new(ExecuteMsg)")
	g.P("err := msg.UnmarshalJSON(messageBytes)")
	g.P("if err != nil { return nil, err }")
	g.P("switch {")
	for _, exec := range g.execs {
		g.P("case msg.", exec.MethodName, " != nil: ")
		g.P("return ", g.typeName, "{}.", exec.MethodName, "(deps, &env, &info, msg.", exec.MethodName, ")")
	}
	g.P("default:")
	g.P("return nil, ", errPkg.Ident("New"), "(\"unknown request\")")
	g.P("}")
	g.P("}")
	g.P()
}

func (g *contractGen) queryFunc() {
	g.P("func Query(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", queryBytes []byte) ([]byte, error) {")
	g.P("query := new(QueryMsg)")
	g.P("err := query.UnmarshalJSON(queryBytes)")
	g.P("if err != nil { return nil, err }")
	g.P("switch {")
	for _, query := range g.queries {
		g.P("case query.", query.MethodName, " != nil:")
		g.P("resp, err := ", g.typeName, "{}.", query.MethodName, "(deps, &env, query.", query.MethodName, ")")
		g.P("if err != nil { return nil, err }")
		g.P("return resp.MarshalJSON()")
	}
	g.P("default:")
	g.P("return nil, ", errPkg.Ident("New"), "(\"unknown request\")")
	g.P("}")
	g.P("}")
	g.P()
}

func (g *contractGen) migrateFunc() {
	g.P("func Migrate(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", messageBytes []byte) (*", typesPkg.Ident("Response"), ", error) {")
	switch g.migrate == nil {
	case true:
		g.P("return &", typesPkg.Ident("Response"), "{}, nil")
	case false:
		var input interface{}
		if g.migrate.Input.ImportPath == "" {
			input = g.migrate.Input.Name
		} else {
			input = protogen.GoIdent{
				GoName:       g.migrate.Input.Name,
				GoImportPath: g.migrate.Input.ImportPath,
			}
		}
		g.P("msg := new(", input, ")")
		g.P("err := msg.UnmarshalJSON(messageBytes)")
		g.P("if err != nil { return nil, err }")
		g.P("return ", g.typeName, "{}.", g.migrate.Name, "(deps, &env, msg)")
	}
	g.P("}")
	g.P()
}

func (g *contractGen) instantiateFunc() {
	g.P("func Instantiate(deps *", stdPkg.Ident("Deps"), ", env ", typesPkg.Ident("Env"), ", info ", typesPkg.Ident("MessageInfo"), ", messageBytes []byte) (*", typesPkg.Ident("Response"), ", error) {")
	switch g.instantiate == nil {
	case true:
		g.P("return &", typesPkg.Ident("Response"), "{}, nil")
	case false:
		var input interface{}
		if g.instantiate.Input.ImportPath == "" {
			input = g.instantiate.Input.Name
		} else {
			input = protogen.GoIdent{
				GoName:       g.instantiate.Input.Name,
				GoImportPath: g.instantiate.Input.ImportPath,
			}
		}
		g.P("msg := new(", input, ")")
		g.P("err := msg.UnmarshalJSON(messageBytes)")
		g.P("if err != nil { return nil, err }")
		g.P("return ", g.typeName, "{}.", g.instantiate.Name, "(deps, &env, &info, msg)")
	}
	g.P("}")
	g.P()
}
