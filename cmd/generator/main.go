package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use: "generator is used to generate boilerplate code for cosmwasm-go contracts.",
	}
	root.AddCommand(stateCmd())
	return root
}

func stateCmd() *cobra.Command {
	state := &cobra.Command{
		Use: "state is used to generates boilerplate for contract's state object",
	}
	state.AddCommand(mapCmd())
	return state
}

func main() {
	root := rootCmd()
	err := root.Execute()
	if err != nil {
		_, _ = fmt.Printf("%s", err)
		os.Exit(1)
	}
}
