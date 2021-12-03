package main

import (
	"fmt"
	"github.com/cosmwasm/cosmwasm-go/pkg/generator"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	generatorPkg = protogen.GoImportPath("github.com/cosmwasm/cosmwasm-go/pkg/generator")
	testingPkg   = protogen.GoImportPath("testing")
)

func mapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "map [file] [type]",
		RunE: func(cmd *cobra.Command, args []string) error {
			xt := filepath.Ext(args[0])
			filename := args[0][:len(args[0])-len(xt)] // strip extension from filename
			genCodePath := fmt.Sprintf("%s.%s.go", filename, strcase.ToSnake(args[1]))
			builderPath := fmt.Sprintf("%s.%s.%s_test.go", filename, strcase.ToSnake(args[1]), "builder")
			testName := fmt.Sprintf("Test_Gen%s", args[1])
			err := generateBuilder(args[0], args[1], builderPath, genCodePath, testName)
			if err != nil {
				return err
			}
			defer os.Remove(builderPath)
			// run codegen test
			goTestPkg := filepath.Dir(args[0])
			switch goTestPkg {
			case ".":
				goTestPkg = "./..."
			default:
				goTestPkg = "./" + goTestPkg
			}
			ex := exec.Command("go", "test", "-run", testName, goTestPkg)
			ex.Stdout = os.Stdout
			ex.Stderr = os.Stderr
			return ex.Run()
		},
		Args: cobra.ExactArgs(2),
	}
	return cmd
}

func generateBuilder(file, typeName, builderPath, genCodePath, testName string) error {
	// we create a test file which creates the type.
	pkgName, err := getPkgName(file)
	if err != nil {
		return err
	}
	g := generator.NewGenerator()
	g.Import(generatorPkg)
	g.Import(testingPkg)
	g.P("// if you find this then delete it")
	g.P("package ", pkgName)

	g.P("func ", testName, "(t *", testingPkg.Ident("T"), ") {")
	g.P("gen, err :=", generatorPkg.Ident("MapState"), "(\"", pkgName, "\", ", typeName, "{})")
	g.P("if err != nil { t.Fatal(err) }")
	g.P("if err := gen.Generate(); err != nil { t.Fatal(err) }")
	g.P("if err := gen.WriteTo(\"", genCodePath, "\"); err != nil { t.Fatal(err) }")
	g.P("}")

	err = g.WriteTo(builderPath)
	if err != nil {
		return err

	}
	return nil
}

func getPkgName(fileName string) (string, error) {
	fs := token.NewFileSet()
	n, err := parser.ParseFile(fs, fileName, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return n.Name.String(), nil
}
