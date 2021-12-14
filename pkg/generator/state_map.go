package generator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	"reflect"
	"strings"
)

const (
	ORMTag        = "orm"
	PrimaryKeyTag = "primaryKey"
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

type mapState struct {
	*Generator // note we use protogen because it takes care of import resolution even with colliding names

	typ        reflect.Type
	primaryKey reflect.StructField

	typeName          string
	stateVar          string
	stateHandler      string
	typeNameSnakeCase string
	pkg               string
}

func (g *mapState) Generate() (err error) {
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
	return nil
}

func (g *mapState) genErrors() {
	g.Import(errPkg)
	g.P("var (")
	g.P("// Err", g.typ.Name(), "NotFound is returned when the object is not found.")
	g.P("Err", g.typ.Name(), "NotFound = ", errPkg.Ident("New"), "(\"", g.pkg, ": ", g.typ.Name(), " not found\")")
	g.P("// Err", g.typ.Name(), "AlreadyExists is returned when the object already exists.")
	g.P("Err", g.typ.Name(), "AlreadyExists = ", errPkg.Ident("New"), "(\"", g.pkg, ": ", g.typ.Name(), " already exists\")")
	g.P(")")
	g.P()
}

func (g *mapState) genPackage() {
	g.P("package ", g.pkg)
}

func (g *mapState) genStateHandlerVar() {
	g.Import(storagePkg)
	g.P("var (")
	g.P("// ", g.stateVar, " provides access to the ", g.typeName, " namespaced storage.")
	g.P(g.stateVar, " = ", g.stateHandler, "{ ns: ", storagePkg.Ident("NewNamespaced"), "(\"", g.typeNameSnakeCase, "\") }")
	g.P(")")
	g.P()
}

func (g *mapState) genStateHandler() {
	g.P("// ", g.stateHandler, " provides ORM functionality around ", g.typeName, ".")
	g.P("type ", g.stateHandler, " struct {")
	g.P("ns ", storagePkg.Ident("Namespaced"))
	g.P("}")
}

func (g *mapState) genCreate() {
	primaryKeyGoType := reflectTypeToGoType(g.primaryKey.Type)

	g.P("// Create handles creation of ", g.typeName, " objects.")
	g.P("// Returns Err", g.typeName, "AlreadyExists in case the object exists.")

	g.P("func (x ", g.stateHandler, ") Create(storage ", stdPkg.Ident("Storage"), ", o ", g.typeName, ") error {")
	// convert field to bytes
	g.P("_k := ", keysPkg.Ident(strcase.ToCamel(primaryKeyGoType+"PrimaryKey")), "(o.", g.primaryKey.Name, ")")
	// check if object exists
	g.P("_v := x.ns.Get(storage, _k)")
	g.P("if _v != nil {")
	g.P("return Err", g.typ.Name(), "AlreadyExists")
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

func (g *mapState) genRead() {
	primaryKeyField := lowerCamelCase(g.primaryKey.Name)
	primaryKeyGoType := reflectTypeToGoType(g.primaryKey.Type)
	g.P("// Read returns ", g.typeName, " given its ", g.primaryKey.Name, ".")
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

func (g *mapState) genUpdate() {
	primaryKeyField := lowerCamelCase(g.primaryKey.Name)
	primaryKeyGoType := reflectTypeToGoType(g.primaryKey.Type)
	g.P("// Update updates an instance of ", g.typeName, ", given its ", g.primaryKey.Name, " by running the provided function f.")
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

func (g *mapState) genDelete() {
	primaryKeyGoType := reflectTypeToGoType(g.primaryKey.Type)
	g.P("// Delete deletes an instance of ", g.typeName, " given its ", g.primaryKey.Name, ".")
	g.P("// Returns Err", g.typeName, "NotFound in case no record is found.")
	g.P("func (x ", g.stateHandler, ") Delete(storage ", stdPkg.Ident("Storage"), ", o ", g.typeName, ") error {")
	g.P("_k := ", keysPkg.Ident(strcase.ToCamel(primaryKeyGoType+"PrimaryKey")), "(o.", g.primaryKey.Name, ")")
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

func MapState(pkg string, o interface{}) (*mapState, error) {
	g := NewGenerator()

	typ := reflect.TypeOf(o)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid object type for %T %s, want struct", o, typ.Kind())
	}
	if pkg == "" {
		s := strings.Split(typ.PkgPath(), "/")
		pkg = s[len(s)-1]
	}

	// find primary key
	var foundFd *reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		// anonymous fields are skipped. TODO(fdymylja): is it good to have anonymous fields inside state objects?
		if fd.Anonymous {
			continue
		}
		v, ok := fd.Tag.Lookup(ORMTag)
		if !ok {
			continue
		}
		if !strings.Contains(v, PrimaryKeyTag) {
			continue
		}
		foundFd = &fd
		break
	}

	if foundFd == nil {
		return nil, fmt.Errorf("no primary key found for %T", o)
	}

	return &mapState{
		Generator:         g,
		typ:               typ,
		typeName:          typ.Name(),
		typeNameSnakeCase: strcase.ToSnake(typ.Name()),
		stateVar:          fmt.Sprintf("%s%s", typ.Name(), stateVarSuffix),
		stateHandler:      fmt.Sprintf("%s%s", typ.Name(), stateHandlerVarSuffix),
		pkg:               pkg,
		primaryKey:        *foundFd,
	}, nil
}
