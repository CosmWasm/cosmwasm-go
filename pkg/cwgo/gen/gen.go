package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
)

type Gen struct {
	*protogen.GeneratedFile
	path string
}

func NewGen(dir string, fileName string) *Gen {
	return &Gen{
		(&protogen.Plugin{}).NewGeneratedFile("xd.go", "un"), filepath.Join(dir, fileName),
	}
}

func (g *Gen) Generate() error {
	c, err := g.Content()
	if err != nil {
		return fmt.Errorf("invalid generated code: %w", err)
	}
	err = os.RemoveAll(g.path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(g.path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(c)
	if err != nil {
		return err
	}

	return nil
}
