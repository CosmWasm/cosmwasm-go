package cwgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	require.NoError(t, Generate("../../example/identityv2/src"))
}
