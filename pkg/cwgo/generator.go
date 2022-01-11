package cwgo

import (
	"fmt"
	"log"

	"github.com/iancoleman/strcase"
)

func Generate(dir string) error {
	p, err := Parse(dir)
	if err != nil {
		return fmt.Errorf("unable to parse dir %s: %w", dir, err)
	}
	if len(p) == 0 {
		log.Printf("nothing to generate in dir %s", dir)
		return nil
	}
	for pkg, parsed := range p {
		gen := newGenerator(dir, pkg, parsed)
		err := gen.Generate()
		if err != nil {
			return fmt.Errorf("unable to generate pkg %s: %w", pkg, err)
		}
	}

	return nil
}

type generator struct {
	p        *Parsed
	dir, pkg string
}

func newGenerator(dir, pkg string, parsed *Parsed) *generator {
	return &generator{
		p:   parsed,
		dir: dir,
		pkg: pkg,
	}
}

func (g *generator) Generate() error {
	// TODO(fdymylja): can we support multiple contracts in same pkg? we could but queryMsg and executeMsg would look ugly.
	if len(g.p.Contracts) > 1 {
		return fmt.Errorf("maximum number of contracts per package supported is 1, got %d", len(g.p.Contracts))
	}
	// first we generate the contract QueryMSG, ExecuteMsg
	for name, contract := range g.p.Contracts {
		gen := newUnionGen(g.dir, g.pkg, contract)
		err := gen.generate()
		if err != nil {
			return fmt.Errorf("unable to generate union types for contract %s: %w", name, err)
		}

		g.p.JSON["QueryMsg"] = &JSON{Name: "QueryMsg"}
		g.p.JSON["ExecuteMsg"] = &JSON{Name: "ExecuteMsg"}
	}
	gen := newJSONGen(g.dir, g.pkg, g.p.JSON)
	err := gen.generate()
	if err != nil {
		return fmt.Errorf("unable to generate json: %w", err)
	}

	for _, o := range g.p.StateObjects {
		gen := newStateMapGen(g.dir, g.pkg, o)
		err := gen.generate()
		if err != nil {
			return fmt.Errorf("unable to generate state object %s: %w", o.Name, err)
		}
	}

	for _, contract := range g.p.Contracts {
		gen := newContractGen(g.dir, g.pkg, contract)
		err := gen.generate()
		if err != nil {
			return fmt.Errorf("unable to generate contract code %s: %w", contract.Name, err)
		}
	}

	return nil
}

func contractFileName(pkg string, name string) string {
	return fmt.Sprintf("zzz_%s_api_%s.go", strcase.ToSnake(pkg), strcase.ToSnake(name))
}

func jsonFileName(pkg string) string {
	return fmt.Sprintf("zzz_%s_json.go", strcase.ToSnake(pkg))
}

func stateObjectFileName(pkg string, typeName string) string {
	return fmt.Sprintf("zzz_%s_state_%s.go", strcase.ToSnake(pkg), strcase.ToSnake(typeName))
}
