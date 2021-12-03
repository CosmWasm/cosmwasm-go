package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
	"os"
)

func NewGenerator() *Generator {
	gen := (&protogen.Plugin{}).NewGeneratedFile("xd.go", "un")
	return &Generator{GeneratedFile: gen}
}

type Generator struct {
	*protogen.GeneratedFile
}

func (g *Generator) WriteTo(path string) error {
	b, err := g.Content()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
