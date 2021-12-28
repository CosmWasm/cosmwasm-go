package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

func contractCmd() *cobra.Command {
	contract := &cobra.Command{
		Use: "contract [file] [type] generates the contract's boilerplate code",
		RunE: func(cmd *cobra.Command, args []string) error {
			xt := filepath.Ext(args[0])
			filename := args[0][:len(args[0])-len(xt)] // strip extension from filename
			genCodePath := fmt.Sprintf("%s.%s.go", filename, "boilerplate")
			builderPath := fmt.Sprintf("%s.%s.%s_test.go", filename, strcase.ToSnake(args[1]), "builder")
			testName := fmt.Sprintf("Test_Gen%s", args[1])
			err := generateBuilder(contract, args[0], args[1], builderPath, genCodePath, testName)
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
			err = ex.Run()
			if err != nil {
				return err
			}
			return nil
		},
		Args: cobra.ExactArgs(2),
	}

	return contract
}
