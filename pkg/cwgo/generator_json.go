package cwgo

import (
	"path/filepath"

	"github.com/CosmWasm/tinyjson/bootstrap"
	"github.com/CosmWasm/tinyjson/parser"
)

func newJSONGen(dir string, pkg string, json map[string]*JSON) *jsonGen {
	return &jsonGen{
		pkg:   pkg,
		dir:   dir,
		types: json,
	}
}

type jsonGen struct {
	pkg, dir string
	types    map[string]*JSON
}

func (g *jsonGen) generate() error {
	jsonTypes := make([]string, 0, len(g.types))
	tinyparser := parser.Parser{
		PkgName: g.pkg,
	}

	for _, t := range g.types {
		jsonTypes = append(jsonTypes, t.Name)
	}

	tinyparser.StructNames = jsonTypes

	err := tinyparser.Parse(g.dir, true)
	if err != nil {
		return err
	}

	tinygen := bootstrap.Generator{
		PkgPath:                  tinyparser.PkgPath,
		PkgName:                  g.pkg,
		Types:                    jsonTypes,
		SnakeCase:                true,
		OmitEmpty:                true,
		DisallowUnknownFields:    true,
		SkipMemberNameUnescaping: true,
		OutName:                  filepath.Join(g.dir, jsonFileName(g.pkg)),
	}

	return tinygen.Run()
}
