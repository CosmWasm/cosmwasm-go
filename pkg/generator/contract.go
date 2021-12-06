package generator

import (
	"reflect"
	"strings"
)

const (
	QueryHandlerPrefix = "Query"
	ExecHandlerPrefix  = "Exec"
)

type ExecDescriptor struct {
}

type QueryDescriptor struct {
}

// Contract takes care of generating the contract boilerplate code.
type Contract struct {
	*Generator
	typ   reflect.Type
	exec  map[string]ExecDescriptor
	query map[string]QueryDescriptor
}

func (c *Contract) Generate() error {
	// we extract the exec and query handlers
	// of the type
	for i := 0; i < c.typ.NumMethod(); i++ {
		method := c.typ.Method(i)
		// we ignore unexported methods.
		if !method.IsExported() {
			continue
		}

		switch {
		case isQuery(method.Name):
			err := c.addQuery(method)
			if err != nil {
				return err
			}
		case isExec(method.Name):
			err := c.addExec(method)
			if err != nil {
				return err
			}
		default: // TODO(fdymylja): do we want to block other things?

		}
	}
	// TODO(fdymylja): embedding

	return nil
}

func (c *Contract) addQuery(method reflect.Method) error {
	funcType := method.Func.Type()
	if funcType.NumIn() != 3 {

	}
	inputs := make([]reflect.Type, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		inputs[i] = funcType.In(i)
	}

	_ = make([]reflect.Type, 0, funcType.NumOut())
	panic("")
}

func (c *Contract) addExec(method reflect.Method) error {
	return nil
}

func isQuery(name string) bool {
	return strings.HasPrefix(name, QueryHandlerPrefix)
}

func isExec(name string) bool {
	return strings.HasPrefix(name, ExecHandlerPrefix)
}
