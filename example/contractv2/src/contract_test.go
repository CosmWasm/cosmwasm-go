package src

import (
	"github.com/cosmwasm/cosmwasm-go/pkg/generator"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGenerateContract(t *testing.T) {
	gen, err := generator.NewContract("src", Contract{})
	require.NoError(t, err)

	err = gen.Generate()
	require.NoError(t, err)

	f, err := os.OpenFile("contract.boilerplate.go", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	require.NoError(t, err)

	defer f.Close()

	require.NoError(t, gen.Write(f))
}
