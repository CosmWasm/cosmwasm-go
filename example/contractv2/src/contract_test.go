package src

import (
	"github.com/cosmwasm/cosmwasm-go/pkg/generator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateContract(t *testing.T) {
	gen, err := generator.NewContract("src", Contract{})
	require.NoError(t, err)

	err = gen.Generate()
	require.NoError(t, err)

	require.NoError(t, gen.WriteTo("contract.boilerplate.go"))
}
