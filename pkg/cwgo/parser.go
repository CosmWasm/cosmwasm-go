package cwgo

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// TODO(fdymylja): we need to type check contract exec, query, migrate, instantiate param types (besides messages)

type StateObjectType int

const (
	StateObjectUndefined StateObjectType = iota
	StateObjectMap
	StateObjectSingleton
)

const (
	GenIdentifier             = "+cw"
	GenSeparator              = ":"
	GenStateIdentifier        = "state"
	GenStateArgumentMap       = "map"
	GenStateArgumentSingleton = "singleton"
	GenExecIdentifier         = "exec"
	GenQueryIdentifier        = "query"
	GenMigrateIdentifier      = "migrate"
	GenInstantiateIdentifier  = "instantiate"
	GenJSONIdentifier         = "json"
)

const (
	OrmPrimaryKeyTag = "orm:\"primaryKey\""
)

const (
	statePrefix       = GenIdentifier + GenSeparator + GenStateIdentifier
	jsonPrefix        = GenIdentifier + GenSeparator + GenJSONIdentifier
	queryPrefix       = GenIdentifier + GenSeparator + GenQueryIdentifier
	execPrefix        = GenIdentifier + GenSeparator + GenExecIdentifier
	migratePrefix     = GenIdentifier + GenSeparator + GenMigrateIdentifier
	instantiatePrefix = GenIdentifier + GenSeparator + GenInstantiateIdentifier
)

type Migrate struct {
	Name  string
	Input *Type
}

type Instantiate struct {
	Name  string
	Input *Type
}

type Contract struct {
	Name        string
	Execs       map[string]*Exec
	Queries     map[string]*Query
	Migrate     *Migrate
	Instantiate *Instantiate
}

type Exec struct {
	MethodName string
	Input      *Type
}

type Query struct {
	MethodName string
	Input      *Type
	Output     *Type
}

type JSON struct {
	Name string
}

type PrimaryKey struct {
	FieldName string
	Type      *Type
}

type Type struct {
	Name       string
	ImportPath protogen.GoImportPath
}

type SecondaryKey struct {
}

type StateObject struct {
	Type          StateObjectType
	Name          string
	PrimaryKey    *PrimaryKey
	SecondaryKeys map[string]*SecondaryKey
}

type Parsed struct {
	astPkg       *ast.Package
	Pkg          string
	Contracts    map[string]*Contract
	JSON         map[string]*JSON
	StateObjects map[string]*StateObject
}

func (p *Parsed) parse(name string, comments *doc.Package) error {
	// set pkg name
	p.Pkg = name

	// find types to generate
	for _, typ := range comments.Types {
		typeName := typ.Name
		if isStateObject(typ.Doc) {
			err := p.addStateObject(typ)
			if err != nil {
				return fmt.Errorf("unable to generate %s as state object: %w", typeName, err)
			}
		}

		if isJSON(typ.Doc) {
			err := p.addJSON(typ)
			if err != nil {
				return fmt.Errorf("unable to generate %s json code", typeName)
			}
		}
	}

	// add contract execs, migrate, query, instantiate
	for _, typ := range comments.Types {
		for _, method := range typ.Methods {
			if isExec(method) {
				err := p.addExec(typ, method)
				if err != nil {
					return fmt.Errorf("unable to add exec method %s of contract %s: %w", typ.Name, method.Name, err)
				}
			}
			if isQuery(method) {
				err := p.addQuery(typ, method)
				if err != nil {
					return fmt.Errorf("unable to add query method %s of contract %s: %w", typ.Name, method.Name, err)
				}
			}

			if isMigrate(method) {
				err := p.addMigrate(typ, method)
				if err != nil {
					return fmt.Errorf("unable to add migrate method %s of contract %s: %w", typ.Name, method.Name, err)
				}
			}

			if isInstantiate(method) {
				err := p.addInstantiate(typ, method)
				if err != nil {
					return fmt.Errorf("unable to add instantiate method %s of contract %s: %w", typ.Name, method.Name, err)
				}
			}
		}
	}
	return nil
}

func (p *Parsed) addStateObject(typ *doc.Type) error {
	stateType := StateObjectUndefined

	for _, line := range strings.Split(typ.Doc, "\n") {
		if !strings.HasPrefix(line, statePrefix) {
			continue
		}

		stateTypeStr := strings.TrimPrefix(line, statePrefix+" ") // we also remove the whitespace
		switch stateTypeStr {
		case GenStateArgumentMap:
			stateType = StateObjectMap
		case GenStateArgumentSingleton:
			stateType = StateObjectSingleton
		default:
			return fmt.Errorf("unrecognized state object type %s", stateTypeStr)
		}
	}

	if stateType == StateObjectUndefined {
		return fmt.Errorf("undefined state object type")
	}

	if len(typ.Decl.Specs) > 1 {
		return fmt.Errorf("unable to handle from ast perspective")
	}

	p.StateObjects[typ.Name] = &StateObject{
		Type:          stateType,
		Name:          typ.Name,
		PrimaryKey:    nil,
		SecondaryKeys: nil,
	}

	if stateType != StateObjectMap {
		return nil
	}

	structType := typ.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)

	var primaryKey *PrimaryKey
	var secondaryKeys map[string]*SecondaryKey // TODO
	for _, field := range structType.Fields.List {
		if field.Tag == nil {
			continue
		}
		if !isPrimaryKey(field.Tag.Value) {
			continue
		}

		fieldType, ok := field.Type.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unable to parse field %#v, note type aliases are not supported", field.Type)
		}
		if len(field.Names) > 1 {
			return fmt.Errorf("unable to parse ast")
		}

		primaryKey = &PrimaryKey{
			FieldName: field.Names[0].Name,
			Type: &Type{
				Name: fieldType.Name,
			},
		}
	}

	p.StateObjects[typ.Name].PrimaryKey = primaryKey
	p.StateObjects[typ.Name].SecondaryKeys = secondaryKeys
	return nil
}

func Parse(dir string) (map[string]*Parsed, error) {
	set := token.NewFileSet()

	d, err := parser.ParseDir(set, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	parsedPkgs := make(map[string]*Parsed, len(d))
	for _, pkg := range d {
		parsed := newParsed(pkg)
		parsed.Pkg = pkg.Name

		comments := doc.New(pkg, "", doc.AllDecls)

		err = parsed.parse(pkg.Name, comments)
		if err != nil {
			return nil, err
		}

		parsedPkgs[pkg.Name] = parsed
	}

	return parsedPkgs, nil
}

func newParsed(astPkg *ast.Package) *Parsed {
	return &Parsed{
		astPkg:       astPkg,
		Pkg:          astPkg.Name,
		Contracts:    map[string]*Contract{},
		JSON:         map[string]*JSON{},
		StateObjects: map[string]*StateObject{},
	}
}

func (p *Parsed) addJSON(typ *doc.Type) error {
	p.JSON[typ.Name] = &JSON{Name: typ.Name}
	return nil
}

func (p *Parsed) addExec(typ *doc.Type, method *doc.Func) error {
	c := p.addContract(typ)

	// TODO(fdymylja): check first params are correct.
	paramType, err := p.methodParamType(typ.Name, 3, method)
	if err != nil {
		return fmt.Errorf("unable to get msg parameter type in contract %s method %s: %w", typ.Name, method.Name, err)
	}

	c.Execs[method.Name] = &Exec{
		MethodName: method.Name,
		Input:      paramType,
	}

	return nil
}

func (p *Parsed) addMigrate(typ *doc.Type, method *doc.Func) error {
	c := p.addContract(typ)
	if c.Migrate != nil {
		return fmt.Errorf("migrate declared twice in contract %s: %s<->%s", typ.Name, c.Migrate.Name, method.Name)
	}

	inputType, err := p.methodParamType(typ.Name, 2, method)
	if err != nil {
		return fmt.Errorf("unable to find input type for migrate method %s in contract %s: %w", method.Name, typ.Name, err)
	}

	c.Migrate = &Migrate{
		Name:  method.Name,
		Input: inputType,
	}

	return nil
}

func (p *Parsed) addInstantiate(typ *doc.Type, method *doc.Func) error {
	c := p.addContract(typ)
	if c.Instantiate != nil {
		return fmt.Errorf("instantiate declared twice in contract %s: %s<->%s", typ.Name, c.Instantiate.Name, method.Name)
	}

	inputType, err := p.methodParamType(typ.Name, 3, method)
	if err != nil {
		return fmt.Errorf("unable to identify input type for instantiation method %s in contract %s: %w", typ.Name, method.Name, err)
	}
	c.Instantiate = &Instantiate{
		Name:  method.Name,
		Input: inputType,
	}

	return nil
}

func (p *Parsed) addQuery(typ *doc.Type, method *doc.Func) error {
	c := p.addContract(typ)
	inputType, err := p.methodParamType(typ.Name, 2, method)
	if err != nil {
		return fmt.Errorf("unable to identify query input type in contract %s method %s: %w", typ.Name, method.Name, err)
	}

	outputType, err := p.methodReturnType(typ.Name, 0, method)
	if err != nil {
		return fmt.Errorf("unable to identify query output type in contract %s method %s: %w", typ.Name, method.Name, err)
	}
	c.Queries[method.Name] = &Query{
		MethodName: method.Name,
		Input:      inputType,
		Output:     outputType,
	}

	return nil
}

func (p *Parsed) addContract(typ *doc.Type) *Contract {
	if contract, exists := p.Contracts[typ.Name]; exists {
		return contract
	}

	p.Contracts[typ.Name] = &Contract{
		Name:        typ.Name,
		Execs:       map[string]*Exec{},
		Queries:     map[string]*Query{},
		Migrate:     nil,
		Instantiate: nil,
	}

	return p.Contracts[typ.Name]
}

func (p *Parsed) methodParamType(structName string, paramPosition int, method *doc.Func) (*Type, error) {
	if method.Decl.Type.Params == nil {
		return nil, fmt.Errorf("function has no params")
	}
	if len(method.Decl.Type.Params.List) < paramPosition+1 {
		return nil, fmt.Errorf("not enough parameters, expected %d, got %d", paramPosition+1, len(method.Decl.Type.Params.List))
	}
	// TODO(fdymylja): check if the first three params are correct
	param := method.Decl.Type.Params.List[paramPosition]
	starExpr := param.Names[0].Obj.Decl.(*ast.Field).Type.(*ast.StarExpr)
	switch m := starExpr.X.(type) {
	// local msg
	case *ast.Ident:
		return &Type{
			Name:       m.Name,
			ImportPath: "",
		}, nil
	// imported msg
	case *ast.SelectorExpr:
		pkg := m.X.(*ast.Ident).Name
		typename := m.Sel.Name

		path, err := p.importPathForMethodParam(structName, method.Name, pkg)
		if err != nil {
			return nil, fmt.Errorf("unable to find import path for package %s in struct's %s method %s: %w", pkg, typename, method.Name, err)
		}

		return &Type{
			Name:       typename,
			ImportPath: protogen.GoImportPath(path),
		}, nil
	default:
		return nil, fmt.Errorf("unable to get method param type: %T", m)
	}
}

func (p *Parsed) importPathForMethodParam(typename string, method string, pkg string) (string, error) {
	for _, f := range p.astPkg.Files {
		for _, decl := range f.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// check if name matches method
			if funcDecl.Name.Name != method {
				continue
			}
			// check if it belongs to struct
			if len(funcDecl.Recv.List) == 0 {
				continue
			}
			// we found our type
			if funcDecl.Recv.List[0].Type.(*ast.Ident).Name != typename {
				continue
			}

			for _, imp := range f.Imports {
				switch {
				case imp.Name != nil:
					if imp.Name.Name == pkg {
						return strings.ReplaceAll(imp.Path.Value, "\"", ""), nil
					}
				default:
					v := strings.ReplaceAll(imp.Path.Value, "\"", "")
					p := strings.Split(v, "/")
					last := p[len(p)-1]
					if last == pkg {
						return imp.Path.Value, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("cannot find method declaration in pkg")
}

func (p *Parsed) methodReturnType(structName string, position int, method *doc.Func) (*Type, error) {
	if method.Decl.Type.Results == nil {
		return nil, fmt.Errorf("function has no returns")
	}

	if len(method.Decl.Type.Results.List) < position+1 {
		return nil, fmt.Errorf("function does not have enough parameters")
	}

	output := method.Decl.Type.Results.List[position].Type.(*ast.StarExpr)

	switch m := output.X.(type) {
	// local msg
	case *ast.Ident:
		return &Type{
			Name:       m.Name,
			ImportPath: "",
		}, nil
	// imported msg
	case *ast.SelectorExpr:
		pkg := m.X.(*ast.Ident).Name
		typename := m.Sel.Name

		path, err := p.importPathForMethodParam(structName, method.Name, pkg)
		if err != nil {
			return nil, fmt.Errorf("unable to find import path for package %s in struct's %s method %s: %w", pkg, typename, method.Name, err)
		}

		return &Type{
			Name:       typename,
			ImportPath: protogen.GoImportPath(path),
		}, nil
	default:
		return nil, fmt.Errorf("unable to get method param type: %T", m)
	}

}

func isJSON(d string) bool {
	return commentHasPrefix(d, jsonPrefix)
}

func isStateObject(d string) bool {
	return commentHasPrefix(d, statePrefix)
}

func isExec(method *doc.Func) bool {
	return commentHasPrefix(method.Doc, execPrefix)
}

func isInstantiate(method *doc.Func) bool {
	return commentHasPrefix(method.Doc, instantiatePrefix)
}

func isMigrate(method *doc.Func) bool {
	return commentHasPrefix(method.Doc, migratePrefix)
}

func isQuery(method *doc.Func) bool {
	return commentHasPrefix(method.Doc, queryPrefix)
}

func commentHasPrefix(doc string, prefix string) bool {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}

	return false
}

func isPrimaryKey(value string) bool {
	return strings.Contains(value, OrmPrimaryKeyTag)
}
