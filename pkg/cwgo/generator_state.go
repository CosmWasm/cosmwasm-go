package cwgo

import (
	"fmt"
	"github.com/cosmwasm/cosmwasm-go/pkg/cwgo/gen"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	errPkg     = protogen.GoImportPath("errors")
	storagePkg = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/std/storage")
	keysPkg    = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/std/storage/keys")
	stdPkg     = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/std")
)

const (
	stateVarSuffix        = "State"
	stateHandlerVarSuffix = stateVarSuffix + "Handler"
)

func newStateMapGen(dir, pkg string, o *StateObject) *stateMapGen {
	return &stateMapGen{
		Gen:           gen.NewGen(dir, stateObjectFileName(pkg, o.Name)),
		pkg:           pkg,
		dir:           dir,
		stateVar:      fmt.Sprintf("%s%s", o.Name, stateVarSuffix),
		stateHandler:  fmt.Sprintf("%s%s", o.Name, stateHandlerVarSuffix),
		typeName:      o.Name,
		typeNameSnake: strcase.ToSnake(o.Name),
		primaryKey:    o.PrimaryKey,
	}
}

type stateMapGen struct {
	*gen.Gen

	pkg, dir      string
	stateVar      string
	stateHandler  string
	typeName      string
	typeNameSnake string
	primaryKey    *PrimaryKey
}

func (g *stateMapGen) generate() (err error) {
	// just panic inside codegen
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	g.genPackage()
	// generate vars and types
	g.genErrors()
	g.genStateHandlerVar()
	// gen state handler
	g.genStateHandler()
	g.genCreate()
	g.genRead()
	g.genUpdate()
	g.genDelete()

	return g.Generate()
}

func (g *stateMapGen) genErrors() {
	g.Import(errPkg)
	g.P("var (")
	g.P("// Err", g.typeName, "NotFound is returned when the object is not found.")
	g.P("Err", g.typeName, "NotFound = ", errPkg.Ident("New"), "(\"", g.pkg, ": ", g.typeName, " not found\")")
	g.P("// Err", g.typeName, "AlreadyExists is returned when the object already exists.")
	g.P("Err", g.typeName, "AlreadyExists = ", errPkg.Ident("New"), "(\"", g.pkg, ": ", g.typeName, " already exists\")")
	g.P(")")
	g.P()
}

func (g *stateMapGen) genPackage() {
	g.P("package ", g.pkg)
}

func (g *stateMapGen) genStateHandlerVar() {
	g.Import(storagePkg)
	g.P("var (")
	g.P("// ", g.stateVar, " provides access to the ", g.typeName, " namespaced storage.")
	g.P(g.stateVar, " = ", g.stateHandler, "{ ns: ", storagePkg.Ident("NewNamespaced"), "(\"", g.typeNameSnake, "\") }")
	g.P(")")
	g.P()
}

func (g *stateMapGen) genStateHandler() {
	g.P("// ", g.stateHandler, " provides ORM functionality around ", g.typeName, ".")
	g.P("type ", g.stateHandler, " struct {")
	g.P("ns ", storagePkg.Ident("Namespaced"))
	g.P("}")
}

func (g *stateMapGen) genCreate() {
	primaryKeyGoType := g.primaryKey.Type.Name

	g.P("// Create handles creation of ", g.typeName, " objects.")
	g.P("// Returns Err", g.typeName, "AlreadyExists in case the object exists.")

	g.P("func (x ", g.stateHandler, ") Create(storage ", stdPkg.Ident("Storage"), ", o ", g.typeName, ") error {")
	// convert field to bytes
	g.P("_k := ", keysPkg.Ident(strcase.ToCamel(primaryKeyGoType+"PrimaryKey")), "(o.", g.primaryKey.FieldName, ")")
	// check if object exists
	g.P("_v := x.ns.Get(storage, _k)")
	g.P("if _v != nil {")
	g.P("return Err", g.typeName, "AlreadyExists")
	g.P("}")
	// marshal object
	g.P("_b, err := o.MarshalJSON()")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P("x.ns.Set(storage, _k, _b)")
	g.P("return nil")
	g.P("}")
	g.P()
}

func (g *stateMapGen) genRead() {
	primaryKeyField := lowerCamelCase(g.primaryKey.Type.Name)
	primaryKeyGoType := g.primaryKey.Type.Name
	g.P("// Read returns ", g.typeName, " given its ", g.primaryKey.FieldName, ".")
	g.P("// Returns Err", g.typeName, "NotFound in case nothing is found.")
	g.P("func (x ", g.stateHandler, ") Read(storage ", stdPkg.Ident("Storage"), ", ", primaryKeyField, " ", primaryKeyGoType, ") (", g.typeName, ", error) {")
	g.P("_k := ", keysPkg.Ident(strcase.ToCamel(primaryKeyGoType+"PrimaryKey")), "(", primaryKeyField, ")")
	g.P("_v := x.ns.Get(storage, _k)")
	g.P("if _v == nil {")
	g.P("return ", g.typeName, "{}, Err", g.typeName, "NotFound")
	g.P("}")
	g.P("_o := new(", g.typeName, ")")
	g.P("err := _o.UnmarshalJSON(_v)")
	g.P("if err != nil {")
	g.P("return ", g.typeName, "{}, err")
	g.P("}")
	g.P("return *_o, nil")
	g.P("}")
}

func (g *stateMapGen) genUpdate() {
	primaryKeyField := lowerCamelCase(g.primaryKey.Type.Name)
	primaryKeyGoType := g.primaryKey.Type.Name
	g.P("// Update updates an instance of ", g.typeName, ", given its ", g.primaryKey.FieldName, " by running the provided function f.")
	g.P("// If f passes a nil ", g.typeName, " it means the object was not found.")
	g.P("// Returning a nil ", g.typeName, " from f will cause no updates in the object.")
	g.P("func (x ", g.stateHandler, ") Update(storage ", stdPkg.Ident("Storage"), ", ", primaryKeyField, " ", primaryKeyGoType, ", f func(old *", g.typeName, ") (*", g.typeName, ", error)) (err error) {")
	g.P("_k := ", keysPkg.Ident(strcase.ToCamel(primaryKeyGoType+"PrimaryKey")), "(", primaryKeyField, ")")
	g.P("_v := x.ns.Get(storage, _k)")
	g.P()
	g.P("var _old *", g.typeName)
	g.P()
	g.P("if _v != nil {")
	g.P("_old = new(", g.typeName, ")")
	g.P("err = _old.UnmarshalJSON(_v)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P("}")
	g.P()
	g.P("_updated, err := f(_old)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	g.P("if _updated == nil {")
	g.P("return nil")
	g.P("}")
	g.P()
	g.P("_v, err = _updated.MarshalJSON()")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	g.P("x.ns.Set(storage, _k, _v)")
	g.P()
	g.P("return nil")
	g.P("}")
}

func (g *stateMapGen) genDelete() {
	primaryKeyGoType := g.primaryKey.Type.Name
	g.P("// Delete deletes an instance of ", g.typeName, " given its ", g.primaryKey.FieldName, ".")
	g.P("// Returns Err", g.typeName, "NotFound in case no record is found.")
	g.P("func (x ", g.stateHandler, ") Delete(storage ", stdPkg.Ident("Storage"), ", o ", g.typeName, ") error {")
	g.P("_k := ", keysPkg.Ident(strcase.ToCamel(primaryKeyGoType+"PrimaryKey")), "(o.", g.primaryKey.FieldName, ")")
	g.P("_v := x.ns.Get(storage, _k)")
	g.P()
	g.P("if _v == nil {")
	g.P("return Err", g.typeName, "NotFound")
	g.P("}")
	g.P()
	g.P("x.ns.Remove(storage, _k)")
	g.P()
	g.P("return nil")
	g.P("}")
}

func lowerCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}
