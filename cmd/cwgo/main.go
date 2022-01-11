package main

import (
	"github.com/cosmwasm/cosmwasm-go/pkg/cwgo"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use: "cwgo [path]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cwgo.Generate(args[0])
		},
		Args: cobra.ExactArgs(1),
	}

	err := cmd.Execute()
	if err != nil {
		cmd.Println(err.Error())
		os.Exit(1)
	}

}
