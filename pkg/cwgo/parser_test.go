package cwgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	res, err := Parse("../../example/identityv2/src")
	require.NoError(t, err)

	j, err := json.Marshal(res)
	require.NoError(t, err)
	t.Logf("%s", j)
}
