package cwgo

import (
	"github.com/cosmwasm/cosmwasm-go/pkg/cwgo/gen"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

func newUnionGen(dir, pkg string, contract *Contract) *unionGen {
	return &unionGen{
		pkg:      pkg,
		dir:      dir,
		contract: contract,
		Gen:      gen.NewGen(dir, contractFileName(pkg, contract.Name)),
	}
}

type unionGen struct {
	pkg, dir string ``
	contract *Contract
	*gen.Gen
}

func (g *unionGen) generate() error {
	g.P("package ", g.pkg)
	g.P("type QueryMsg struct {")
	for name, query := range g.contract.Queries {
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
	for name, exec := range g.contract.Execs {
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
	return g.Generate()
}
